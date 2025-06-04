let sum = BigInt(0);
const start = performance.now();
for (let i = 0; i < 1_000_000_000; i++) {
  sum += BigInt(i);
}
const end = performance.now();
console.log(`Node Loop Benchmark: ${(end - start).toFixed(3)}ms`);
console.log(`Sum: ${sum}`);
