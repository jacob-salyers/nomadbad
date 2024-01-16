import { promises as fs } from 'fs';
import { randomUUID } from 'crypto';

export class Student {
    constructor(obj) {
        this.id = obj.id;
        this.first_name = obj.first_name;
        this.last_name = obj.last_name;
    }
}

Student.new = (f, l) => new Student(generateUUID(), f, l);
Student.read = () =>
	fs.readFile('data/students.json', { encoding: 'utf8' })
        .then(JSON.parse)
        .then(arr => arr.map(o => new Student(o)));

Student.write = (students) => 
    fs.writeFile('data/students.json', JSON.stringify(students));

export class Class {
    constructor(obj) {
        this.id = obj.id;
        this.days = obj.days;
        this.display_title = obj.display_title;
    }
}
