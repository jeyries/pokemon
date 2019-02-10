
extern crate rayon;

use rayon::prelude::*;
use rayon::iter::{IntoParallelIterator, ParallelIterator};
use std::env;
use std::fs::File;
use std::io::BufReader;
use std::io::BufRead;
use std::time::Duration;
use std::time::SystemTime;

const M: usize = 150;
const N: usize = 6;

//static mut _names: Option<Vec<String>> = None ;
//static mut _attacks: [[u64; 3]; M] = [[0; 3]; M];

fn to_float_secs(d: Duration) -> f64 {
    return d.as_secs() as f64 + d.subsec_nanos() as f64 * 1e-9;
}

struct State {
    names: Vec<String>,
    attacks: [[u64; 3]; M],
}

struct Result {
    best: u32,
    answer: [u8; N],
}

impl State {

     #[inline(never)]
    fn search_i0(& self, result: &mut Result, i0: usize) {

        let mut s4: [u64; 3] = [0; 3];
        let attacks = & self.attacks;

        for i1 in 0 .. i0 {
            for i2 in 0 .. i1 {
                for i3 in 0 .. i2 {
                    for i4 in 0 .. i3 {

                        s4[0] = 0;
                        s4[1] = 0;
                        s4[2] = 0;

                        s4[0] |= attacks[i0][0];
                        s4[1] |= attacks[i0][1];
                        s4[2] |= attacks[i0][2];

                        s4[0] |= attacks[i1][0];
                        s4[1] |= attacks[i1][1];
                        s4[2] |= attacks[i1][2];

                        s4[0] |= attacks[i2][0];
                        s4[1] |= attacks[i2][1];
                        s4[2] |= attacks[i2][2];

                        s4[0] |= attacks[i3][0];
                        s4[1] |= attacks[i3][1];
                        s4[2] |= attacks[i3][2];

                        s4[0] |= attacks[i4][0];
                        s4[1] |= attacks[i4][1];
                        s4[2] |= attacks[i4][2];

                        for i5 in 0 .. i4 {

                            let mut s50 = s4[0];
                            let mut s51 = s4[1];
                            let mut s52 = s4[2];

                            s50 |= attacks[i5][0];
                            s51 |= attacks[i5][1];
                            s52 |= attacks[i5][2];

                            let mut score: u32 = s50.count_ones();
                            score += s51.count_ones();
                            score += s52.count_ones();

                            if score > result.best {
                                result.best = score;
                                result.answer[0] = i0 as u8;
                                result.answer[1] = i1 as u8;
                                result.answer[2] = i2 as u8;
                                result.answer[3] = i3 as u8;
                                result.answer[4] = i4 as u8;
                                result.answer[5] = i5 as u8;
                            }

                        }

                    }
                }
            }
        }
    }

    fn map(& self, i0: usize) -> Result {

        println!("[i0 {}] processing", i0);
        let start = SystemTime::now();

        let mut result = Result {best: 0, answer: [0; N]};
        self.search_i0(&mut result, i0);

        let elapsed = start.elapsed().unwrap();
        println!("[i0 {}] elapsed = {:.3} s", i0, to_float_secs(elapsed));
        println!("[i0 {}] best = {}", i0, result.best);

        return result
    }

    fn reduce(& self, _current: Result, result: Result) -> Result {

        let mut current = _current;

        if result.best > current.best {
            current.best = result.best;
            for i in 0 .. N {
                current.answer[i] = result.answer[i];
            }
        }

        println!("[current] best = {}", current.best);
        
        return current
    }

    fn search(& self) -> Result {

        return (0 .. M).rev()
            .map(|i0| self.map(i0) )
            .fold( Result {best: 0, answer: [0; N]}, 
                |current, result| self.reduce(current, result));

    }

    fn par_search(& self) -> Result {

        return (0 .. M).into_par_iter()
            .map(|i0| self.map(i0) )
            .reduce(|| Result {best: 0, answer: [0; N]}, 
                    |current, result| self.reduce(current, result));  
        

    }

    fn load(&mut self, f: File) {

        let names = &mut self.names;
        let attacks = &mut self.attacks;

        let mut k: usize = 0;
        let file = BufReader::new(&f);
        for line in file.lines() {
            let l = line.unwrap();

            let mut split = l.split(",");

            let name = split.next().unwrap();
            names.push(name.to_string());

            for s in split {
                let x: usize = s.parse::<usize>().unwrap();
                attacks[k][x/64] |= 1 << (x % 64)
            }

            k += 1;
        }    

        println!("{} {:016x} {:016x} {:016x}", names[0], attacks[0][0], attacks[0][1], attacks[0][2]);
    }
}





fn main() {
    let args: Vec<String> = env::args().collect();

    if args.len() < 2 {
        println!("usage: pokemon input.txt [--parallel]");
        return
	}

    let mut filename: Option<String> = None;
    let mut parallel = false;

    for arg in args {
        if arg == "--parallel" {
            parallel = true;
        } else {
            filename = Some(arg); 
        }
    }

    let mut state = State { names: Vec::new(), attacks: [[0; 3]; M] };

    let f = File::open(filename.unwrap()).expect("file not found");
    state.load(f);

    let start = SystemTime::now();

    let result: Result;
    if parallel {
        result = state.par_search();
    } else {
        result = state.search();
    }
 
    let elapsed = start.elapsed().unwrap();
    println!("elapsed = {:.3} s", to_float_secs(elapsed));
	println!("best = {}", result.best);

	for i in 0 .. N {
		println!("{}", state.names[result.answer[i] as usize])
	}

}

