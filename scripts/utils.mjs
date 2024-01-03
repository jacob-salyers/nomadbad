import { Student } from './types.mjs';


async function newStudent(f, l) {
    Student.read()
        .then(students => {
            students.push(new Student(f, l));
            return students;
        })
        .then(students => Student.write(students));
}

const funcs = {
    'new-student': (arr) => newStudent(arr.shift(), arr.shift())
};

const args = process.argv.slice(2);
funcs[args.shift()](args);
