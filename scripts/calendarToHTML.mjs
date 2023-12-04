
import { promises as fs } from 'fs';
import { parseICS, generateHTML } from './calendarToHTML-header.mjs';

const schedulePath = process.argv[2] ?? '/dev/stdout';
const eventsPath = process.argv[3];

process.stdin.setEncoding('utf8');

new Promise((resolve, reject) => {
	let acc  = '';
	process.stdin.on('data', data => acc += data);
	process.stdin.on('error', err => reject(err));
	process.stdin.on('end', () => resolve(acc));
})
	.then(parseICS)
    .then(generateHTML)
	.then(events => {
		const promises = [];
		promises.push(fs.writeFile(schedulePath, events.recurring));
		if (eventsPath != null)
			promises.push(fs.writeFile(eventsPath, events.upcoming));

		return Promise.all(promises);
	});
