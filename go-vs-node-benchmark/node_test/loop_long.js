const Long = require("long");
let sum = Long.fromInt(0);
const start = performance.now();
for (let i = 0; i < 1_000_000_000; i++) {
  sum = sum.add(Long.fromInt(i));
}
const end = performance.now();
console.log(`Node Loop Benchmark: ${(end - start).toFixed(3)}ms`);
console.log(`Sum: ${sum.toString()}`);
