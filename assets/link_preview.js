// @feature:layouts Client-side link preview on hover for internal links (desktop only).
(function () {
  var BASE_URL = "{{.BaseURL}}";
  var HOVER_DELAY = 300;
  var LEAVE_GRACE = 100;
  var TOOLTIP_MAX_WIDTH = 400;
  var TOOLTIP_MAX_HEIGHT = 300;

  var cache = new Map();
  var hoverTimer = null;
  var leaveTimer = null;
  var initialized = false;

  function getTooltip() {
    var el = document.getElementById("link-preview-tooltip");
    if (el) return el;

    el = document.createElement("div");
    el.id = "link-preview-tooltip";
    el.style.position = "absolute";
    el.style.maxWidth = TOOLTIP_MAX_WIDTH + "px";
    el.style.maxHeight = TOOLTIP_MAX_HEIGHT + "px";
    el.style.overflowY = "auto";
    el.style.display = "none";
    el.style.zIndex = "9999";
    document.body.appendChild(el);

    el.addEventListener("mouseenter", function () {
      clearTimeout(leaveTimer);
    });
    el.addEventListener("mouseleave", function () {
      hideTooltip();
    });

    return el;
  }

  function hideTooltip() {
    var el = document.getElementById("link-preview-tooltip");
    if (el) el.style.display = "none";
  }

  function positionTooltip(tooltip, anchor) {
    var rect = anchor.getBoundingClientRect();
    var scrollX = window.pageXOffset || document.documentElement.scrollLeft;
    var scrollY = window.pageYOffset || document.documentElement.scrollTop;
    var spaceBelow = window.innerHeight - rect.bottom;
    var top;

    if (spaceBelow >= TOOLTIP_MAX_HEIGHT + 8) {
      top = rect.bottom + scrollY + 8;
    } else {
      top = rect.top + scrollY - TOOLTIP_MAX_HEIGHT - 8;
    }

    var left = rect.left + scrollX;
    if (left + TOOLTIP_MAX_WIDTH > window.innerWidth + scrollX) {
      left = window.innerWidth + scrollX - TOOLTIP_MAX_WIDTH - 8;
    }
    if (left < scrollX) {
      left = scrollX + 8;
    }

    tooltip.style.top = top + "px";
    tooltip.style.left = left + "px";
  }

  function shouldSkipLink(link) {
    if (link.getAttribute("target") === "_blank") return true;

    var href = link.getAttribute("href");
    if (!href) return true;
    if (href.charAt(0) === "#") return true;

    var currentPath = window.location.pathname;
    try {
      var linkUrl = new URL(href, window.location.origin);
      if (linkUrl.pathname === currentPath) return true;
    } catch (e) {
      return true;
    }

    return false;
  }

  function fetchPreview(href) {
    if (cache.has(href)) return Promise.resolve(cache.get(href));

    return fetch(href)
      .then(function (res) {
        if (!res.ok) throw new Error("fetch failed");
        return res.text();
      })
      .then(function (html) {
        var parser = new DOMParser();
        var doc = parser.parseFromString(html, "text/html");
        var contentEl = doc.querySelector("#content .content");
        var result = contentEl ? contentEl.innerHTML : "";
        cache.set(href, result);
        return result;
      })
      .catch(function () {
        cache.set(href, "");
        return "";
      });
  }

  function onLinkEnter(e) {
    var link = e.currentTarget;
    if (shouldSkipLink(link)) return;

    clearTimeout(leaveTimer);
    clearTimeout(hoverTimer);

    hoverTimer = setTimeout(function () {
      var href = link.getAttribute("href");
      fetchPreview(href).then(function (content) {
        if (!content) return;
        var tooltip = getTooltip();
        tooltip.innerHTML = content;
        positionTooltip(tooltip, link);
        tooltip.style.display = "block";
      });
    }, HOVER_DELAY);
  }

  function onLinkLeave() {
    clearTimeout(hoverTimer);
    leaveTimer = setTimeout(function () {
      hideTooltip();
    }, LEAVE_GRACE);
  }

  function bindLinks() {
    var links = document.querySelectorAll("a.internal-link");
    for (var i = 0; i < links.length; i++) {
      if (links[i].dataset.previewBound) continue;
      links[i].dataset.previewBound = "1";
      links[i].addEventListener("mouseenter", onLinkEnter);
      links[i].addEventListener("mouseleave", onLinkLeave);
    }
  }

  function onScroll() {
    hideTooltip();
  }

  function onKeydown(e) {
    if (e.key === "Escape") hideTooltip();
  }

  function onClickOutside(e) {
    var tooltip = document.getElementById("link-preview-tooltip");
    if (!tooltip) return;
    if (tooltip.style.display === "none") return;
    if (tooltip.contains(e.target)) return;
    if (e.target.closest && e.target.closest("a.internal-link")) return;
    hideTooltip();
  }

  window.initLinkPreview = function () {
    if (initialized) {
      bindLinks();
      return;
    }

    if (window.innerWidth < 1024) return;
    if (!window.matchMedia("(hover: hover)").matches) return;

    initialized = true;

    getTooltip();
    bindLinks();

    window.addEventListener("scroll", onScroll, true);
    document.addEventListener("keydown", onKeydown);
    document.addEventListener("click", onClickOutside);
  };
})();
