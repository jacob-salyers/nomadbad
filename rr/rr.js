import fs from 'fs';
import { Player } from './glicko.js';

export function demo_init() {
	if (fs.existsSync('competitors.json'))
		return Tournament.fromJSON('competitors.json');
	else 
		return new Tournament([
			new Player('Jacob'),
			new Player('DuBois'),
			new Player('Marc'),
			new Player('Nowell'),
			new Player('Emily'),
			new Player('Spencer'),
			new Player('Max')
		]);
}

export function demo_persist(t) {
	t.toJSON('competitors.json');
}

export class Tournament {

	constructor(competitors = [], match = 1, matches = []) {
		this.competitors = competitors;
		this.match = match;
		this.matches = matches;
	}

	add_new(name) {
		this.competitors.push(new Player(name));
	}

	add(player) {
		this.competiors.push(player);
	}

	generate_matches() {
		// Randomize competitors - Only relevant if players have the 
		// same rating
		this.competitors.sort((a,b) => .5 - Math.random());

		/*
		 * Triangle sort - Sort by rating then rebuild the list so 
		 * that highest rated players are in the middle.
		 * This has the impact of favoring them for Byes
		 */
		this.competitors.sort((a,b) => b.r - a.r);
		const tmp = [];
		let toggle = true;
		for (const el of this.competitors) {
			if (toggle) 
				tmp.unshift(el);
			else 
				tmp.push(el);

			toggle = !toggle;
		}
		this.competitors = tmp;

		// Ensure even number of competitors by adding dummy
		if (this.competitors.length % 2 !== 0)
			this.competitors.push(null);

		/*
		 * Generate the matches
		 */
		for (let round = 0; round < this.competitors.length - 1; round++) {

			/*
			 * Generate a round of matches
			 */

			let i = 0, j = this.competitors.length - 1;

			while (i < j ) {
				const match = [ 
					this.competitors[i++],
					this.competitors[j--]
				];

				// Don't add matches with dummy competitor
				if (!match.includes(null))
					this.matches.push(match);
			}

			/*
			 * Rotate the competitor list
			 */
			const last = this.competitors.pop();
			this.competitors.splice(1,0,last);
		}

		// Remove dummy competitor if exists
		this.competitors = this.competitors.filter(c => c != null);
	}

	report_result(i, res) {
		const match = this.matches[i-1];
		if (match == null) {
			console.log(`Match number ${i} does not exist.`);
			return;
		}

		if (match.length > 2) {
			console.log('This match result has already been finalized.');
			return;
		}
			
		switch (res) {
			case 'left':
				match.push(match[0].name);
				match[0].reportWin(match[1]);
				match[1].reportLoss(match[0]);
				break;
			case 'right':
				match.push(match[1].name);
				match[1].reportWin(match[0]);
				match[0].reportLoss(match[1]);
				break;
			case 'draw':
				match.push('Draw');
				match[1].reportDraw(match[0]);
				match[0].reportDraw(match[1]);
				break;
		}
	}

	left() {
		this.report_result(this.match++, 'left');
		console.clear();
		this.print_matches();
	}

	right() {
		this.report_result(this.match++, 'right');
		console.clear();
		this.print_matches();
	}

	draw() {
		this.report_result(this.match++, 'draw');
		console.clear();
		this.print_matches();
	}

	finish() {
		this.competitors.forEach(c => c.updateRating());
	}

	print_players(i = 0) {
		const maxlen = this.competitors.reduce((m,c) => (c?.name.length ?? 0)  > m ? c.name.length : m, 0);
		for (; i < this.competitors.length; i++) {
			const name = this.competitors[i]?.name ?? 'Bye';
			let padding = '    ';
			for (let j = 0; j < maxlen - name.length; j++)
				padding += ' ';

			const rating = this.competitors[i]?.r ?? 'N/A';
			console.log(`${name}:${padding}${rating}`);
		}
	}

	print_matches() {
		const c1 = [];
		const c2 = [];
		const c3 = [];

		let c1max = 0, c2max = 0;
		this.matches.forEach((match, i) => {
			const str1 = `Match ${i+1}:`;
			if (str1.length > c1max)
				c1max = str1.length;

			c1.push(str1);

			const str2 = `${match[0].name} vs ${match[1].name}`;
			if (str2.length > c2max)
				c2max = str2.length;

			c2.push(str2);

			if (match[2] != null)
				c3.push(`(${match[2]})`);
			else 
				c3.push('');
		});

		for (let i = 0; i < this.matches.length; i++) {
			let p1 = '    ';
			for (let j = 0; j < c1max - c1[i].length; j++)
				p1 += ' ';

			let p2 = '   ';
			for (let j = 0; j < c2max - c2[i].length; j++)
				p2 += ' ';

			console.log(c1[i]+p1+c2[i]+p2+c3[i]);
		}
	}

	toJSON(path) {
		fs.writeFileSync(
			path,
			JSON.stringify(this.competitors),
			{ encoding: 'utf8' }
		);
	}

	static fromJSON(path) {
		const competitors = JSON.parse(fs.readFileSync(path, { encoding: 'utf8' }));
		return new Tournament(competitors.map(Player.fromJSON));
	}
}
