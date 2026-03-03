// @feature:search Client-side full-text search with inverted index and dropdown results.
(function () {
  var BASE_URL = "{{.BaseURL}}";
  var MAX_RESULTS = 15;
  var SNIPPET_LEN = 120;
  var MIN_QUERY_LEN = 2;

  var invertedIndex = null;
  var indexEntries = null;

  function fetchIndex() {
    if (window._searchIndex) {
      return Promise.resolve(window._searchIndex);
    }
    return fetch(BASE_URL + "/search-index.json")
      .then(function (res) { return res.json(); })
      .then(function (data) {
        window._searchIndex = data;
        return data;
      });
  }

  function buildInvertedIndex(entries) {
    var idx = {};
    for (var i = 0; i < entries.length; i++) {
      var text = ((entries[i].title || "") + " " + (entries[i].content || "")).toLowerCase();
      var words = text.split(/\s+/);
      for (var w = 0; w < words.length; w++) {
        var word = words[w].replace(/[^a-z0-9]/g, "");
        if (!word) continue;
        if (!idx[word]) idx[word] = new Set();
        idx[word].add(i);
      }
    }
    return idx;
  }

  function searchEntries(query) {
    if (!invertedIndex || !indexEntries) return [];
    var tokens = query.toLowerCase().split(/\s+/).filter(function (t) {
      return t.replace(/[^a-z0-9]/g, "").length > 0;
    });
    if (tokens.length === 0) return [];

    var resultSet = null;
    var indexKeys = Object.keys(invertedIndex);

    for (var t = 0; t < tokens.length; t++) {
      var token = tokens[t].replace(/[^a-z0-9]/g, "");
      if (!token) continue;
      var matching = new Set();
      for (var k = 0; k < indexKeys.length; k++) {
        if (indexKeys[k].indexOf(token) === 0) {
          var entries = invertedIndex[indexKeys[k]];
          entries.forEach(function (idx) { matching.add(idx); });
        }
      }
      if (resultSet === null) {
        resultSet = matching;
      } else {
        var intersection = new Set();
        resultSet.forEach(function (idx) {
          if (matching.has(idx)) intersection.add(idx);
        });
        resultSet = intersection;
      }
    }

    if (!resultSet) return [];
    var results = [];
    resultSet.forEach(function (idx) { results.push(indexEntries[idx]); });
    return results.slice(0, MAX_RESULTS);
  }

  function snippet(content, query) {
    if (!content) return "";
    var lower = content.toLowerCase();
    var tokens = query.toLowerCase().split(/\s+/);
    var pos = -1;
    for (var i = 0; i < tokens.length; i++) {
      pos = lower.indexOf(tokens[i]);
      if (pos >= 0) break;
    }
    var start = Math.max(0, pos - 30);
    var s = content.substring(start, start + SNIPPET_LEN);
    if (start > 0) s = "..." + s;
    if (start + SNIPPET_LEN < content.length) s = s + "...";
    return s;
  }

  function getOrCreateDropdown(input) {
    var existing = document.getElementById("search-results");
    if (existing) return existing;
    var dropdown = document.createElement("div");
    dropdown.id = "search-results";
    dropdown.style.position = "absolute";
    dropdown.style.zIndex = "9999";
    input.parentNode.style.position = "relative";
    input.parentNode.appendChild(dropdown);
    return dropdown;
  }

  function hideDropdown() {
    var dd = document.getElementById("search-results");
    if (dd) dd.style.display = "none";
  }

  function showResults(input, results, query) {
    var dropdown = getOrCreateDropdown(input);
    dropdown.innerHTML = "";
    if (results.length === 0) {
      dropdown.style.display = "none";
      return;
    }
    dropdown.style.display = "";
    for (var i = 0; i < results.length; i++) {
      var entry = results[i];
      var item = document.createElement("a");
      item.href = entry.url;
      item.className = "search-result-item";
      item.setAttribute("data-index", i);

      var title = document.createElement("div");
      title.className = "search-result-title";
      title.textContent = entry.title || "";
      item.appendChild(title);

      if (entry.folder) {
        var folder = document.createElement("div");
        folder.className = "search-result-folder";
        folder.textContent = entry.folder;
        item.appendChild(folder);
      }

      var snip = document.createElement("div");
      snip.className = "search-result-snippet";
      snip.textContent = snippet(entry.content, query);
      item.appendChild(snip);

      dropdown.appendChild(item);
    }
  }

  function filterSidebar(term) {
    var items = document.querySelectorAll("#left-sidebar li");
    items.forEach(function (item) {
      var text = item.textContent.toLowerCase();
      var matches = text.includes(term);
      item.style.display = matches ? "" : "none";
      if (matches && term) {
        var details = item.querySelector("details");
        if (details) details.open = true;
        var parent = item.parentElement;
        while (parent && parent.closest(".sidebar")) {
          if (parent.tagName === "DETAILS") parent.open = true;
          parent = parent.parentElement;
        }
      }
    });
  }

  function getHighlightedIndex(dropdown) {
    var items = dropdown.querySelectorAll(".search-result-item");
    for (var i = 0; i < items.length; i++) {
      if (items[i].classList.contains("highlighted")) return i;
    }
    return -1;
  }

  function setHighlighted(dropdown, index) {
    var items = dropdown.querySelectorAll(".search-result-item");
    for (var i = 0; i < items.length; i++) {
      items[i].classList.remove("highlighted");
    }
    if (index >= 0 && index < items.length) {
      items[index].classList.add("highlighted");
      items[index].scrollIntoView({ block: "nearest" });
    }
  }

  window.initFullSearch = function () {
    var searchInput = document.getElementById("navbar-search");
    if (!searchInput) return;

    var newInput = searchInput.cloneNode(true);
    searchInput.parentNode.replaceChild(newInput, searchInput);

    fetchIndex().then(function (data) {
      indexEntries = data;
      invertedIndex = buildInvertedIndex(data);
    });

    newInput.addEventListener("input", function (e) {
      var term = e.target.value.trim();
      if (term.length < MIN_QUERY_LEN) {
        hideDropdown();
        filterSidebar(term.toLowerCase());
        return;
      }
      var results = searchEntries(term);
      showResults(newInput, results, term);
    });

    newInput.addEventListener("keydown", function (e) {
      var dropdown = document.getElementById("search-results");
      if (!dropdown || dropdown.style.display === "none") return;
      var items = dropdown.querySelectorAll(".search-result-item");
      if (items.length === 0) return;

      var current = getHighlightedIndex(dropdown);

      if (e.key === "ArrowDown") {
        e.preventDefault();
        setHighlighted(dropdown, Math.min(current + 1, items.length - 1));
      } else if (e.key === "ArrowUp") {
        e.preventDefault();
        setHighlighted(dropdown, Math.max(current - 1, 0));
      } else if (e.key === "Enter") {
        e.preventDefault();
        if (current >= 0 && items[current]) {
          window.location.href = items[current].href;
        }
      } else if (e.key === "Escape") {
        hideDropdown();
      }
    });

    document.addEventListener("click", function (e) {
      var dropdown = document.getElementById("search-results");
      if (!dropdown) return;
      if (!dropdown.contains(e.target) && e.target !== newInput) {
        hideDropdown();
      }
    });
  };
})();
