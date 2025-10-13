use std::f64::consts::PI;

pub type Matrix = Vec<Vec<f64>>;

pub struct DctPoint {
    x_max: usize,
    y_max: usize,
    x_scales: [f64; 2],
    y_scales: [f64; 2],
}

impl DctPoint {
    pub fn new(x: usize, y: usize, x_s: [f64; 2], y_s: [f64; 2]) -> DctPoint {
        DctPoint {
            x_max: x,
            y_max: y,
            x_scales: x_s,
            y_scales: y_s,
        }
    }

    pub fn calculate(&self, image_data: &Matrix, u: usize, v: usize) -> f64 {
        let mut sum = 0.0f64;
        for i in 0..self.x_max {
            for j in 0..self.y_max {
                let image_value = image_data[i][j];
                let arg_x = ((2 * i + 1) as f64) * (u as f64) * PI / (2.0 * self.x_max as f64);
                let arg_y = ((2 * j + 1) as f64) * (v as f64) * PI / (2.0 * self.y_max as f64);
                sum += image_value * arg_x.cos() * arg_y.cos();
            }
        }
        sum * self.find_scale_factor(u, v)
    }

    fn find_scale_factor(&self, u: usize, v: usize) -> f64 {
        let x_scale = if u == 0 { self.x_scales[0] } else { self.x_scales[1] };
        let y_scale = if v == 0 { self.y_scales[0] } else { self.y_scales[1] };
        x_scale * y_scale
    }
}

pub fn find_dct_matrix(matrix: Matrix) -> Matrix {
    let x_max = matrix.len();
    let y_max = matrix.get(0).map(|r| r.len()).unwrap_or(0);
    assert!(x_max > 0 && y_max > 0, "empty matrix");

    let x_s = [1.0 / (x_max as f64).sqrt(), (2.0 / x_max as f64).sqrt()];
    let y_s = [1.0 / (y_max as f64).sqrt(), (2.0 / y_max as f64).sqrt()];

    let dct_point = DctPoint::new(x_max, y_max, x_s, y_s);

    let mut dct_matrix = vec![vec![0.0f64; y_max]; x_max];
    for u in 0..x_max {
        for v in 0..y_max {
            dct_matrix[u][v] = dct_point.calculate(&matrix, u, v);
        }
    }
    dct_matrix
}

pub fn reduce_matrix(dct_matrix: &Matrix, size: usize) -> Matrix {
    let x_max = dct_matrix.len();
    let y_max = dct_matrix[0].len();
    assert!(size <= x_max && size <= y_max, "reduce size too big");

    let mut new = vec![vec![0.0f64; size]; size];
    for x in 0..size {
        for y in 0..size {
            new[x][y] = dct_matrix[x][y];
        }
    }
    new
}

pub fn calculate_median_value(dct_matrix: &Matrix) -> f64 {
    let n = dct_matrix.len();
    let m = dct_matrix[0].len();
    let mut vals: Vec<f64> = Vec::with_capacity(n * m);
    for x in 0..n {
        for y in 0..m {
            if x == 0 && y == 0 { continue; }
            vals.push(dct_matrix[x][y]);
        }
    }
    if vals.is_empty() {
        return 0.0;
    }
    vals.sort_by(|a, b| a.partial_cmp(b).unwrap());
    let mid = vals.len() / 2;
    if vals.len() % 2 == 1 {
        vals[mid]
    } else {
        (vals[mid - 1] + vals[mid]) / 2.0
    }
}

pub fn build_hash(dct_matrix: &Matrix, threshold: f64) -> String {
    let x_size = dct_matrix.len();
    let y_size = dct_matrix[0].len();
    let mut hash = String::with_capacity(x_size * y_size);
    for x in 0..x_size {
        for y in 0..y_size {
            if dct_matrix[x][y] > threshold {
                hash.push('1');
            } else {
                hash.push('0');
            }
        }
    }
    hash
}

pub fn find_distance(hash1: &str, hash2: &str) -> Option<usize> {
    if hash1.len() != hash2.len() { return None; }
    Some(hash1.chars().zip(hash2.chars()).filter(|(a,b)| a != b).count())
}
