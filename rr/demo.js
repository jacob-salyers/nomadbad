
let r = await import('./rr.js');
let t = r.demo_init();
t.generate_matches();
console.clear(); t.print_matches();
