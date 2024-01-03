
import { Student, Class } from './types.mjs';
import { generateSignIn } from './calendarToHTML-header.mjs';

Promise.all([
	fs.readFile('data/classes.json', { encoding: 'utf8' })
        .then(JSON.parse)
        .then(arr => arr.map(o => new Class(o))),
	Student.read()
]).then(res => generateSignIn(res[0], res[1]));
