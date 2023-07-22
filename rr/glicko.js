import crypto from 'crypto';

/*
 * https://en.wikipedia.org/wiki/Glicko_rating_system
 */

/********************************************************************
 *							  CONSTANTS                             *
 ********************************************************************/

// Rating periods to max uncertainty
const mu = 100; // TODO: This is a guess
// Average player's ratings deviation
const avg_rd = 50; // TODO: This is a guess

const c = Math.sqrt((mu**2 - avg_rd**2)/100);
// ln(10)/400
const q = 0.00575646273;

/********************************************************************
 *						  HELPER FUNCTIONS		                    *
 ********************************************************************/

const get_diff_in_weeks = (d0, d1) => {
	// Get difference in milliseconds
	const diff = d1.getTime() - d0.getTime();

	// Convert to num days and return
	return Math.floor(diff / (1000 * 3600 * 24 * 7));
}

// g(RD_i)
const g = (RD_i) => 1/Math.sqrt(1 + (3 * q**2 * RD_i**2)/Math.PI**2);
// E(s|r_0, r_i, RD_i)
const E = (r_0, r_i, RD_i) => 1 / 
	(1 + 10**((g(RD_i)*(r_0 - r_i))/-400))

const calc_d2 = (r0, games) => 1/(q**2 * games.reduce((sum, opp) => 
	sum += g(opp.rd)**2 * E(r0, opp.r, opp.rd) * (1 - E(r0, opp.r, opp.rd)),
	0));

/*
 * Primary steps
 */

/** 
 * Step 1: Determine ratings deviation
 */
const RD_add_uncertainty = (rd0, t) => Math.min(
	Math.sqrt(rd0**2 + t*c**2),
	350
);

/**
 * Step 2: Determine new rating
 */
const update_rating = (r0, rd, d2, games) => 
	r0 + (q/(1/rd**2 + 1/d2))*games.reduce((sum, opp) => 
		sum += g(opp.rd)*(opp.s - E(r0, opp.r, opp.rd)),
		0);

/**
 * Step 3: Determine new ratings deviation
 */
const update_rd = (rd0, d2) => Math.sqrt(1/(1/rd0**2 + 1/d2));

/********************************************************************
 *							HELPER CLASS							*
 ********************************************************************/

class Game {
	constructor(s, r, rd) {
		this.s = s;
		this.r = r;
		this.rd = RD_add_uncertainty(rd, new Date());
	}
}

/********************************************************************
 *							 EXPORTS								*
 ********************************************************************/

export class Player {
	constructor(name, id = crypto.randomUUID(), r = 1500, rd = 350, last_date = new Date()) {
		this.name = name;
		this.id = id;
		this.r = r;
		this.rd = rd;
		this.last_date = last_date;
		this.games = [];
	}

	updateRating()  {
		const t = get_diff_in_weeks(this.last_date, new Date());
		const rd0 = RD_add_uncertainty(this.rd, t);
		const d2 = calc_d2(rd0, this.games);

		this.r = update_rating(this.r, rd0, d2, this.games);
		this.rd = update_rd(rd0, d2);
		this.games = [];
		this.last_date = new Date();

	}

	reportWin(opponent) {
		this.games.push(new Game(1, opponent.r, opponent.rd));
	}

	reportLoss(opponent) {
		this.games.push(new Game(0, opponent.r, opponent.rd));
	}
	
	reportDraw(opponent) {
		this.games.push(new Game(0.5, opponent.r, opponent.rd));
	}
}

Player.fromJSON = (json) => 
	new Player(json.name, json.id, json.r, json.rd, new Date(json.last_date));
