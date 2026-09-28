document.getElementById("shorten-form").addEventListener("submit", async (e) => {
  e.preventDefault();
  const url = document.getElementById("url-input").value.trim();
  const alias = document.getElementById("alias-input").value.trim();
  const result = document.getElementById("result");
  const btn = document.getElementById("submit-btn");
  result.innerHTML = "";
  if (!url) {
    showError("Please enter a URL");
    return;
  }
  const apiBase = window.location.protocol.startsWith("http")
    ? ""
    : "http://localhost:8080";
  btn.disabled = true;
  btn.classList.add("loading");
  const originalText = btn.textContent;
  btn.textContent = "Shortening…";
  try {
    const res = await fetch(`${apiBase}/shorten`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ url, alias }),
    });
    if (!res.ok) {
      const text = await res.text();
      showError(`Error ${res.status}: ${text || "could not shorten URL"}`);
      return;
    }
    const data = await res.json();
    const origin = window.location.protocol.startsWith("http")
      ? window.location.origin
      : "http://localhost:8080";
    const shortUrl = `${origin}/${data.short_url}`;
    result.innerHTML = "";
    const card = document.createElement("div");
    card.className = "result-card";
    const link = document.createElement("a");
    link.href = shortUrl;
    link.target = "_blank";
    link.rel = "noopener";
    link.textContent = shortUrl;
    const copy = document.createElement("button");
    copy.type = "button";
    copy.className = "copy-btn";
    copy.textContent = "Copy";
    copy.addEventListener("click", async () => {
      try {
        await navigator.clipboard.writeText(shortUrl);
      } catch (_) {
        const ta = document.createElement("textarea");
        ta.value = shortUrl;
        document.body.appendChild(ta);
        ta.select();
        document.execCommand("copy");
        ta.remove();
      }
      copy.textContent = "Copied ✓";
      copy.classList.add("copied");
      setTimeout(() => {
        copy.textContent = "Copy";
        copy.classList.remove("copied");
      }, 1600);
    });
    card.appendChild(link);
    card.appendChild(copy);
    result.appendChild(card);
  } catch (err) {
    console.error(err);
    showError(
      `Request failed: ${err.message}. Is the server running on http://localhost:8080?`
    );
  } finally {
    btn.disabled = false;
    btn.classList.remove("loading");
    btn.textContent = originalText;
  }
});

function showError(message) {
  const result = document.getElementById("result");
  result.innerHTML = "";
  const card = document.createElement("div");
  card.className = "error-card";
  card.textContent = message;
  result.appendChild(card);
  void card.offsetWidth;
  card.classList.add("shake");
}
