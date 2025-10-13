#![feature(unwrap_infallible)]

use crate::phash::*;
use clap::Parser;
use rayon::prelude::*;
use std::collections::HashMap;
use std::fs::{remove_file, DirEntry};
use std::path::PathBuf;

mod phash;

const MIN_DISTANCE: usize = 10;

#[derive(Debug, Clone, Eq)]
struct Pic {
    name: String,
    hash: String,
    size: u64,
}

impl PartialEq for Pic {
    fn eq(&self, other: &Self) -> bool {
        self.hash == other.hash
    }
}

impl Pic {
    fn find_size(&self) -> String {
        const KB: f64 = 1024.0;
        let sizef = self.size as f64;
        if sizef < KB {
            format!("{:.2}b", sizef)
        } else if sizef < KB.powi(2) {
            format!("{:.2}kb", sizef / KB)
        } else if sizef < KB.powi(3) {
            format!("{:.2}mb", sizef / KB.powi(2))
        } else if sizef < KB.powi(4) {
            format!("{:.2}gb", sizef / KB.powi(3))
        } else {
            format!("{}b", self.size)
        }
    }
}

#[derive(clap::ValueEnum, Clone, PartialEq, Default, Debug)]
enum Recursive {
    #[default]
    Non,
    Segmented,
    Flat,
}

#[derive(Parser)]
#[command(version, about)]
struct Args {
    #[arg(short, long)]
    quiet: bool,

    #[arg(long)]
    rm: bool,

    #[arg(short, long, default_value_t, value_enum)]
    recursive: Recursive,

    #[arg(short, long, default_value_t = 0)]
    deep: u32,

    #[arg(short, long)]
    path: Option<PathBuf>,
}

fn main() {
    let args = Args::parse();
    let path = args.path.clone().unwrap_or_else(|| PathBuf::from("./"));
    check(path, args);
}

fn read_dir(dir: PathBuf) -> Vec<DirEntry> {
    let formats = ["png", "jpg", "jpeg"];
    let mut result = Vec::new();

    let rd = match dir.read_dir() {
        Ok(it) => it,
        Err(_) => return result,
    };

    for entry in rd.flatten() {
        let path = entry.path();
        if let Some(ext_os) = path.extension() {
            if let Some(ext_str) = ext_os.to_str() {
                let ext_low = ext_str.to_ascii_lowercase();
                if formats.contains(&ext_low.as_str()) {
                    result.push(entry);
                }
            }
        }
    }
    result
}
fn mk_file_index(
    res: &mut HashMap<PathBuf, Vec<DirEntry>>,
    root: PathBuf,
    dir: PathBuf,
    is_flat: bool,
    deep: u32,
) {
    // читаем файлы в текущей директории
    let mut files: Vec<DirEntry> = read_dir(dir.clone());

    // вставляем в индекс в зависимости от режима
    if !files.is_empty() {
        if is_flat {
            // flat: собираем ВСЁ под ключом root (перемещаем файлы в бакет root)
            let bucket = res.entry(root.clone()).or_insert_with(Vec::new);
            bucket.append(&mut files); // переносим элементы из files в bucket
        } else {
            // segmented (и non—тоже): добавляем текущую директорию как отдельный сегмент,
            // но только если в ней >= 2 файлов (как у тебя было раньше)
            if files.len() >= 2 {
                res.insert(dir.clone(), files); // перемещаем files в map (без клонирования)
            } // иначе — не вставляем
        }
    }

    // если глубина 0 — не рекурсить дальше
    if deep == 0 {
        return;
    }

    // рекурсивный обход подпапок (без паники)
    let rd = match dir.read_dir() {
        Ok(it) => it,
        Err(_) => return,
    };

    for entry in rd.flatten() {
        let p = entry.path();
        if p.is_dir() {
            mk_file_index(
                res,
                root.clone(),
                p,
                is_flat,
                deep.saturating_sub(1),
            );
        }
    }
}

fn check(dir: PathBuf, args: Args) {
    let mode = args.recursive.clone();
    let mut files: HashMap<PathBuf, Vec<DirEntry>> = HashMap::new();
    mk_file_index(&mut files, dir.clone(), dir.clone(), mode == Recursive::Flat, args.deep);

    if files.is_empty() {
        println!("No directories with >=2 supported images found under {}", dir.display());
        return;
    }

    println!("calculation...");

    let mut pics: HashMap<PathBuf, Vec<Option<Pic>>> = HashMap::new();
    for (k, v) in files.iter() {
        let key = k.clone();
        let vec_of_pics: Vec<Option<Pic>> = v
            .par_iter()
            .map(|e| {
                let path = e.path();
                let name = path.to_string_lossy().to_string();
                if !args.quiet {
                    println!("{}", name);
                }
                let hash = match find_hash(&name) {
                    Some(h) => h,
                    None => return None,
                };
                let size = e.metadata().ok().map(|m| m.len()).unwrap_or(0);
                Some(Pic { name, hash, size })
            })
            .collect();
        pics.insert(key, vec_of_pics);
    }

    pics.par_iter()
        .for_each(|(dir, pics)| process_pics(dir, pics, args.rm))
}

fn process_pics(dir: &PathBuf, pics: &[Option<Pic>], rm: bool) {
    let pics_vec: Vec<Pic> = pics.iter().filter_map(|p| p.clone()).collect();
    let n = pics_vec.len();
    if n < 2 {
        println!("{:?}: no duplicates found (too few files)", dir);
        return;
    }

    let mut pairs: Vec<(usize, usize, usize)> = Vec::new();
    for i in 0..n {
        for j in (i + 1)..n {
            if let Some(d) = find_distance(&pics_vec[i].hash, &pics_vec[j].hash) {
                pairs.push((i, j, d));
            }
        }
    }

    if pairs.is_empty() {
        println!("{:?}: no comparable hashes (length mismatch?)", dir);
        return;
    }

    pairs.sort_by_key(|t| t.2);

    let mut duplicates: Vec<(Pic, Pic)> = Vec::new();
    for (i, j, d) in pairs.iter() {
        if *d < MIN_DISTANCE {
            duplicates.push((pics_vec[*i].clone(), pics_vec[*j].clone()));
        }
    }

    if duplicates.is_empty() {
        println!("{:?}: no duplicates found.", dir);
        return;
    }

    let mut s = String::new();
    for (a, b) in duplicates {
        if rm {
            let d = if a.size < b.size { a.clone() } else { b.clone() };
            remove_file(d.name.clone()).unwrap_or_else(|_| eprintln!("cant delete {}", d.name));
        }
        s += format!(
            "\t{}, {} -- {}, {}\n",
            a.name.clone(),
            a.find_size(),
            b.name.clone(),
            b.find_size()
        )
            .as_str();
    }
    if !s.is_empty() && s.ends_with('\n') {
        s.pop();
    }
    println!("{:?}:\n{}", dir, s);
}

fn find_hash(path: &str) -> Option<String> {
    let size = 50u32;
    let img = match image::open(path) {
        Ok(img) => img
            .resize_to_fill(size, size, image::imageops::Lanczos3)
            .grayscale(),
        Err(_) => return None,
    };

    let gray = img.to_luma8();
    let (x_size, y_size) = gray.dimensions();
    let mut matrix: Matrix = Vec::with_capacity(x_size as usize);
    for x in 0..x_size {
        let mut col = Vec::with_capacity(y_size as usize);
        for y in 0..y_size {
            col.push(gray.get_pixel(x, y)[0] as f64);
        }
        matrix.push(col);
    }

    let dct_matrix = find_dct_matrix(matrix);
    let small = reduce_matrix(&dct_matrix, 8);
    let threshold = calculate_median_value(&small);
    let hash = build_hash(&small, threshold);
    Some(hash)
}
