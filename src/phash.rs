#[allow(dead_code)]
use image::{DynamicImage};
use std::f64::consts::PI;
use std::ops::Add;
use std::panic;
use std::str::Chars;

// MIN_DISTANCE - threshold of duplicates distance
pub const MIN_DISTANCE: i32 = 3;

pub type Matrix = Vec<Vec<f64>>;

pub struct DctPoint<'a> {
    x_max: i64,
    y_max: i64,
    x_scales: &'a mut [f64; 2],
    y_scales: &'a mut [f64; 2],
}

impl DctPoint<'_> {
    pub fn new<'a>(x: i64, y: i64, x_s: &'a mut [f64; 2], y_s: &'a mut [f64; 2]) -> DctPoint<'a> {
        DctPoint {
            x_max: x,
            y_max: y,
            x_scales: x_s,
            y_scales: y_s,
        }
    }

    pub fn calculate(&self, image_data: &Matrix, x: i64, y: i64) -> f64 {
        let mut sum = 0.;
        for i in 0..self.x_max {
            for j in 0..self.y_max {
                let image_value = image_data[i as usize][j as usize];
                let fst_cosine =
                    ((((1 + (2 * i)) * x) as f64) * PI / (2. * self.x_max as f64)).cos();
                let snd_cosine =
                    ((((1 + (2 * j)) * y) as f64) * PI / (2. * self.y_max as f64)).cos();
                let s = image_value * fst_cosine * snd_cosine;
                sum += s;
            }
        }
        sum * self.find_scale_factor(x, y)
    }

    pub fn find_scale_factor(&self, x: i64, y: i64) -> f64 {
        let mut x_scale_factor = self.x_scales[1];
        if x == 0 {
            x_scale_factor = self.x_scales[0];
        }
        let mut y_scale_factor = self.y_scales[1];
        if y == 0 {
            y_scale_factor = self.y_scales[0];
        }
        x_scale_factor * y_scale_factor
    }
}

pub fn find_distance(hash1: &Chars, hash2: &Chars) -> i32 {
    hash1
        .clone()
        .zip(hash2.clone())
        .fold(0, |acc, x| if x.0 != x.1 { acc + 1 } else { acc })
}

pub fn find_hash(img: String) -> Option<String> {
    let size = 50;
    let img = match panic::catch_unwind(|| {
        image::open(img)
            .unwrap()
            .resize_to_fill(size, size, image::imageops::Lanczos3)
            .grayscale()
    }) {
        Ok(img) => img,
        Err(_) => return None,
    };
    let image_matrix = find_image_matrix(img);
    let dct_matrix = find_dct_matrix(image_matrix);
    let small_matrix = reduce_matrix(dct_matrix, 10);
    let dct_mean_value = calculate_mean_value(&small_matrix);
    Some(build_hash(small_matrix, dct_mean_value))
}

fn find_image_matrix(img: DynamicImage) -> Matrix {
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
    matrix
}

pub fn find_dct_matrix(matrix: Matrix) -> Matrix {
    let x_max = matrix.len();
    let y_max = matrix[0].len();
    let mut x_s = [1. / (x_max as f64).sqrt(), (2. / x_max as f64).sqrt()];
    let mut y_s = [1. / (y_max as f64).sqrt(), (2. / y_max as f64).sqrt()];
    let dct_point = DctPoint::new (
        x_max as i64,
        y_max as i64,
        &mut x_s,
        &mut y_s,
    );
    let mut dct_matrix: Matrix = Vec::new();
    for x in 0..x_max {
        dct_matrix.push(Vec::new());
        for y in 0..y_max {
            dct_matrix[x].push(dct_point.calculate(&matrix, x as i64, y as i64));
        }
    }
    dct_matrix
}

pub fn reduce_matrix(dct_matrix: Matrix, size: i64) -> Matrix {
    let mut new_matrix: Matrix = Vec::new();
    for x in 0..size {
        new_matrix.push(Vec::new());
        for y in 0..size {
            new_matrix[x as usize].push(dct_matrix[x as usize][y as usize]);
        }
    }
    new_matrix
}

pub fn calculate_mean_value(dct_matrix: &Matrix) -> f64 {
    let mut avg = 0.;
    let n = dct_matrix.len();
    for x in 0..n {
        for y in (x + 1)..n {
            avg += dct_matrix[x][y] / (n * n) as f64;
        }
    }
    avg
}

pub fn build_hash(dct_matrix: Matrix, dct_mean_value: f64) -> String {
    let mut hash = String::new();
    let x_size = dct_matrix.len();
    let y_size = dct_matrix[0].len();
    for x in 0..x_size {
        for y in 0..y_size {
            if dct_matrix[x][y] > dct_mean_value {
                hash = hash.add("1");
            } else {
                hash = hash.add("0")
            }
        }
    }
    hash
}

#[cfg(test)]
mod phash_tests {
    #[test]
    fn find_distance_test() {
        use super::*;
        assert_eq!(find_distance(&"1101".chars(), &"1011".chars()), 2);
        assert_eq!(find_distance(&"1111".chars(), &"1111".chars()), 0);
    }
}
