import express from "express";
const app = express();
const PORT = 3000;

// Benchmark 1: loop test
function loopTest() {
  let sum = 0;
  for (let i = 1; i <= 1_000_000_000; i++) {
    sum += i;
  }
  console.log("Sum:", sum);
}

// Benchmark 2: concurrent 100ms x 10 tasks
async function concurrencyTest() {
  const wait = () => new Promise((r) => setTimeout(r, 100));
  const tasks = Array.from({ length: 10 }, () => wait());
  await Promise.all(tasks);
}

// HTTP endpoint
app.get("/ping", (req, res) => {
  res.send("pong");
});

app.listen(PORT, () => {
  console.log(`Node.js listening on port ${PORT}`);
  // Uncomment one of the following to test:
  // loopTest();
  // concurrencyTest();
});
