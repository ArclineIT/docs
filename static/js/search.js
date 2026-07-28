(function () {
  'use strict';

  var index = null;
  var input = document.getElementById('search-input');
  var results = document.getElementById('search-results');
  if (!input || !results) return;

  function root() {
    var meta = document.querySelector('meta[name="docs-root"]');
    return meta ? meta.content : '';
  }

  function loadIndex(cb) {
    if (index) { cb(); return; }
    fetch(root() + 'search.json')
      .then(function (r) { return r.json(); })
      .then(function (data) { index = data; cb(); })
      .catch(function () {});
  }

  function esc(s) {
    return s.replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;').replace(/"/g, '&quot;');
  }

  function search(q) {
    if (!index) return;
    q = q.trim();
    if (!q) { hide(); return; }

    var terms = q.toLowerCase().split(/\s+/).filter(Boolean);
    var matches = index.filter(function (doc) {
      var text = (doc.title + ' ' + doc.section + ' ' + doc.excerpt).toLowerCase();
      return terms.every(function (t) { return text.indexOf(t) !== -1; });
    }).slice(0, 8);

    if (!matches.length) { hide(); return; }

    results.innerHTML = matches.map(function (doc) {
      return '<a href="' + esc(doc.url) + '" class="search-result">' +
        '<span class="search-result__section">' + esc(doc.section) + '</span>' +
        '<span class="search-result__title">' + esc(doc.title) + '</span>' +
        '<span class="search-result__excerpt">' + esc(doc.excerpt) + '</span>' +
        '</a>';
    }).join('');
    results.hidden = false;
  }

  function hide() { results.hidden = true; }

  input.addEventListener('focus', function () { loadIndex(function () {}); });

  input.addEventListener('input', function () {
    loadIndex(function () { search(input.value); });
  });

  document.addEventListener('click', function (e) {
    if (!e.target.closest('.search-wrap')) hide();
  });

  document.addEventListener('keydown', function (e) {
    if (e.key === 'Escape') { hide(); input.blur(); }
    if ((e.metaKey || e.ctrlKey) && e.key === 'k') {
      e.preventDefault();
      input.focus();
      input.select();
    }
  });

  // Keyboard navigation inside results
  results.addEventListener('keydown', function (e) {
    var links = results.querySelectorAll('.search-result');
    var active = results.querySelector('.search-result:focus');
    var idx = Array.prototype.indexOf.call(links, active);
    if (e.key === 'ArrowDown') {
      e.preventDefault();
      var next = links[idx + 1] || links[0];
      if (next) next.focus();
    } else if (e.key === 'ArrowUp') {
      e.preventDefault();
      var prev = links[idx - 1] || links[links.length - 1];
      if (prev) prev.focus();
    }
  });

  input.addEventListener('keydown', function (e) {
    if (e.key === 'ArrowDown' && !results.hidden) {
      e.preventDefault();
      var first = results.querySelector('.search-result');
      if (first) first.focus();
    }
  });
})();
