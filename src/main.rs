use rand::Rng;
use std::error::Error;
use std::fmt;
use std::fs::File;
use std::io::{BufRead, BufReader};

// Structures
#[derive(Debug, Clone)]
struct DataPoint {
    features: Vec<f64>,
    label: f64,
}

#[derive(Debug)]
struct Dataset {
    data: Vec<DataPoint>,
    feature_count: usize,
}

// Model traits
trait Model {
    fn predict(&self, features: &[f64]) -> f64;
    fn train(&mut self, dataset: &Dataset) -> Result<(), Box<dyn Error>>;
}

// Linear Regression
struct LinearRegression {
    weights: Vec<f64>,
    bias: f64,
    learning_rate: f64,
    iterations: usize,
}

impl LinearRegression {
    fn new(feature_count: usize, learning_rate: f64, iterations: usize) -> Self {
        let mut rng = rand::thread_rng();
        LinearRegression {
            weights: (0..feature_count).map(|_| rng.gen::<f64>()).collect(),
            bias: rng.gen::<f64>(),
            learning_rate,
            iterations,
        }
    }
}

impl Model for LinearRegression {
    fn train(&mut self, dataset: &Dataset) -> Result<(), Box<dyn Error>> {
        for _ in 0..self.iterations {
            let mut gradient_weights = vec![0.0; self.weights.len()];
            let mut gradient_bias = 0.0;

            for data_point in &dataset.data {
                let prediction = self.predict(&data_point.features);
                let error = prediction - data_point.label;

                for (i, &feature) in data_point.features.iter().enumerate() {
                    gradient_weights[i] += error * feature;
                }
                gradient_bias += error;
            }

            for i in 0..self.weights.len() {
                self.weights[i] -= self.learning_rate * gradient_bias / dataset.data.len() as f64;
            }
        }
        Ok(())
    }

    fn predict(&self, features: &[f64]) -> f64 {
        features
            .iter()
            .zip(&self.weights)
            .map(|(&x, &w)| x * w)
            .sum::<f64>()
            + self.bias
    }
}



fn main() {
    println!("Hello, world!");
}
