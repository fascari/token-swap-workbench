const statusEl = document.querySelector("#health-status");

async function updateHealth() {
  try {
    const response = await fetch("/health", { cache: "no-store" });
    statusEl.textContent = response.ok ? "API online" : "API error";
  } catch {
    statusEl.textContent = "API offline";
  }
}

void updateHealth();
