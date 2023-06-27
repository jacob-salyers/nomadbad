use uuid::Uuid;
use serde::{Serialize,Deserialize};
use serde_json;

use chrono::DateTime;
use chrono::offset::Utc;

/*
 * https://en.wikipedia.org/wiki/Glicko_rating_system
 */

/****************************************************************
 *							  CONSTANTS                         *
 ****************************************************************/

const INIT_RATING: f32 = 1500.0;
const INIT_RD: f32 = 350.0;

// Rating periods to max uncertainty
const MU: f32 = 100.0; // TODO: This is a guess

// Average player's ratings deviation
const AVG_RD: f32 = 50.0; // TODO: This is a guess

// ln(10)/400
const Q: f32 = 0.00575646273;
const PI: f32 = 3.14159;

fn get_diff_in_weeks(d0: DateTime<Utc>, d1: DateTime<Utc>) -> u64 {
    (d1.timestamp() - d0.timestamp()) as u64
}

fn g(rd_i: f32) -> f32 {
    1.0 / (1.0 + (3.0*Q.powf(2.0)*rd_i.powf(2.0))/PI.powf(2.0)).sqrt()
}

fn e(r_0: f32, r_i: f32, rd_i:f32) -> f32 {
    1.0 / (1.0 + 10.0_f32.powf((g(rd_i)*(r_0 - r_i)) / 400.0))
}

fn calc_d2(r0: f32, games: &Vec<Game>) -> f32 {
    1.0 / (Q.powf(2.0) * games.iter()
           .fold(0.0, |sum: f32, o| sum + g(o.rd).powf(2.0)
                 * e(r0, o.r, o.rd)
                 * (1.0 - e(r0, o.r, o.rd))))
}

fn rd_add_uncertainty(rd_0: f32, t: u64) -> f32 {
    let c : f32 = ((MU.powf(2.0) - AVG_RD.powf(2.0)) / 100.0).sqrt();
    (rd_0.powf(2.0) + (t as f32)*c.powf(2.0)).sqrt()
}

fn update_rating(r0: f32, rd: f32, d2: f32, games: &Vec<Game>) -> f32 {
    r0 + (Q/(1.0/rd.powf(2.0) + 1.0/d2))*games.iter()
        .fold(0.0, |sum, o| sum + g(o.rd)*(o.s - e(r0, o.r, o.rd)))
}

fn update_rd(rd_0: f32, d2: f32) -> f32 {
    (1.0 / (1.0 / rd_0.powf(2.0) + 1.0 / d2)).sqrt()
}

#[derive(Serialize, Deserialize)]
struct Game {
    s: f32,
    r: f32,
    rd: f32
}

#[derive(Serialize, Deserialize)]
pub struct Player {
    pub name: String,
    pub id: Uuid,
    pub r: f32,
    rd: f32,
    last_date: DateTime<Utc>,
    games: Vec<Game>
}

impl Player {
    pub fn new(name: String) -> Player {
        Player {
            name,
            id: Uuid::new_v4(),
            r: INIT_RATING,
            rd: INIT_RD,
            last_date: Utc::now(),
            games: Vec::new()
        }
    }

    pub fn update_rating(&mut self) {
        let t = get_diff_in_weeks(self.last_date, Utc::now());
        let rd0 = rd_add_uncertainty(self.rd, t);
        let d2 = calc_d2(rd0, &self.games);

        self.r = update_rating(self.r, rd0, d2, &self.games);
        self.rd = update_rd(rd0, d2);

        self.last_date = Utc::now();
        self.games.clear();
    }

    pub fn report_win(&mut self, opponent: &Player) {
        self.games.push(Game { s: 1.0, r: opponent.r, rd: opponent.rd });
    }

    pub fn report_loss(&mut self, opponent: &Player) {
        self.games.push(Game { s: 0.0, r: opponent.r, rd: opponent.rd });
    }

    pub fn report_draw(&mut self, opponent: &Player) {
        self.games.push(Game { s: 0.5, r: opponent.r, rd: opponent.rd });
    }

    pub fn from_json(s: &str) -> Result<Player, serde_json::Error> {
        serde_json::from_str(s)
    }
}

