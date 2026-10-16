(function () {
  "use strict";

  // Theme toggle, remembered across visits. Storage can be unavailable, so guard it.
  var root = document.documentElement;
  var saved = null;
  try { saved = localStorage.getItem("theme"); } catch (e) { /* ignore */ }
  if (saved) { root.setAttribute("data-theme", saved); }

  var toggle = document.getElementById("theme-toggle");
  if (toggle) {
    toggle.addEventListener("click", function () {
      var next = root.getAttribute("data-theme") === "dark" ? "light" : "dark";
      root.setAttribute("data-theme", next);
      try { localStorage.setItem("theme", next); } catch (e) { /* ignore */ }
    });
  }

  // Routes table: fetch JSON from the server and render it.
  var tbody = document.querySelector("#routes-table tbody");
  if (tbody) {
    fetch("data/routes.json")
      .then(function (res) {
        if (!res.ok) { throw new Error("HTTP " + res.status); }
        return res.json();
      })
      .then(function (routes) {
        routes.forEach(function (route) {
          var tr = document.createElement("tr");
          [route.match, route.action, route.target, route.note].forEach(function (text) {
            var td = document.createElement("td");
            td.textContent = text;
            tr.appendChild(td);
          });
          tbody.appendChild(tr);
        });
      })
      .catch(function (err) {
        tbody.innerHTML = "<tr><td colspan='4'>Could not load routes: " + err.message + "</td></tr>";
      });
  }

  // Latency test: request a small file repeatedly, bypassing the browser cache.
  var ping = document.getElementById("ping");
  var out = document.getElementById("ping-result");
  if (ping && out) {
    ping.addEventListener("click", function () {
      var rounds = 10;
      var times = [];
      ping.disabled = true;
      out.textContent = "Running...";

      function next(n) {
        if (n === rounds) { return Promise.resolve(); }
        var start = performance.now();
        return fetch("img/logo.svg?n=" + n + "&t=" + Date.now(), { cache: "no-store" })
          .then(function (res) { return res.arrayBuffer(); })
          .then(function () {
            times.push(performance.now() - start);
            return next(n + 1);
          });
      }

      next(0)
        .then(function () {
          var sum = times.reduce(function (a, b) { return a + b; }, 0);
          out.textContent = rounds + " requests: min " + Math.min.apply(null, times).toFixed(1) +
            " ms, avg " + (sum / rounds).toFixed(1) + " ms, max " + Math.max.apply(null, times).toFixed(1) + " ms";
        })
        .catch(function (err) { out.textContent = "Failed: " + err.message; })
        .then(function () { ping.disabled = false; });
    });
  }
})();
