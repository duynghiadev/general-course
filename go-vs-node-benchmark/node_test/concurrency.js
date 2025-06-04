async function delay(ms) {
  return new Promise((resolve) => setTimeout(resolve, ms));
}

async function task() {
  await delay(100);
}

async function main() {
  const start = performance.now();
  const tasks = Array(10)
    .fill()
    .map(() => task());
  await Promise.all(tasks);
  const end = performance.now();
  console.log(`Node Concurrency: ${(end - start).toFixed(1)}ms`);
}

main();
