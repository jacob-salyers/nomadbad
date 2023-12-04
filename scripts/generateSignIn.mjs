
import { promises as fs } from 'fs';
import { generateSignIn } from './calendarToHTML-header.mjs';

Promise.all([
	fs.readFile('data/classes.json', { encoding: 'utf8' }).then(JSON.parse),
	fs.readFile('data/students.json', { encoding: 'utf8' }).then(JSON.parse)
]).then(res => generateSignIn(res[0], res[1]));
