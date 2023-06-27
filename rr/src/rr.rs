use rand::thread_rng;
use rand::seq::SliceRandom;
use std::fs;
use std::error::Error;
use serde_json;
use strum;
use crate::glicko::Player;

#[derive(strum::Display)]
pub enum MatchResult {
    Win,
    Lose,
    Draw,
    NoContest
}

pub fn demo_init() -> Result<Tournament<'static>, Box<dyn Error>> {
    let o = fs::read_to_string("competitors.json")?;
    Ok(Tournament {
        competitors: serde_json::from_str(&*o)?,
        count: 1,
        matches: Vec::new()
    })
}

 pub fn demo_persist(t: &Tournament) -> Result<(), Box<dyn Error>> {
     Ok(serde_json::to_writer(
             fs::File::create("competitors.json")?,
             &t.competitors)?)
 }

pub struct Tournament <'a> {
    competitors: Vec<Player>,
    count: u32,
    matches: Vec<(&'a Player, &'a Player)>
}

impl <'a> Tournament <'a> {
    pub fn add_new(&mut self, name: String) {
        self.competitors.push(Player::new(name));
    }

    pub fn add(&mut self, player: Player) {
        self.competitors.push(player);
    }

    pub fn generate_matches(&'a mut self) {

        // randomize
        self.competitors.shuffle(&mut thread_rng());

        /*
         * Triangle sort - Sort by rating then rebuild the list so 
         * that highest rated players are in the middle.
         * self has the impact of favoring them for Byes
         */

        self.competitors.sort_by(|a, b| b.r.partial_cmp(&a.r).unwrap());

        let mut toggle = true;
        let mut tmp: Vec<Option<&Player>> = self.competitors
            .iter()
            .fold(vec![],| mut tmp, p| {
                if toggle {
                    tmp.insert(0, Some(p));
                } else {
                    tmp.push(Some(p));
                }

                toggle = !toggle;
                tmp
            });

        if tmp.len() % 2 != 0 {
            tmp.push(None);
        }

        for _ in 0..(tmp.len() - 1) {
            let (mut i, mut j) = (0, tmp.len() - 1);
            while i < j {
                if let (Some(a), Some(b)) = (tmp[i], tmp[j]) {
                    self.matches.push((a, b));
                }
                i += 1; j += 1;

                let last = tmp.pop().unwrap();
                tmp.insert(1, last);
            }
        }
    }


    report_result(&mut self, i: u32, res: MatchResult) {
    }
}
