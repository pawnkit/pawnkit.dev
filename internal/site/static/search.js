(() => {
  const input = document.querySelector('#search-query');
  if (!input) return;

  const kind = document.querySelector('#search-kind');
  const results = document.querySelector('#search-results');
  const status = document.querySelector('#search-status');
  const params = new URLSearchParams(location.search);
  let entries = [];

  input.value = params.get('q') || '';
  kind.value = params.get('kind') || '';

  const render = () => {
    const query = input.value.trim().toLowerCase();
    const selected = kind.value;
    const matches = entries.filter((entry) => {
      const text = `${entry.title} ${entry.summary} ${entry.source}`.toLowerCase();
      return (!selected || entry.kind === selected) && (!query || text.includes(query));
    }).slice(0, 100);

    const items = matches.map((entry) => {
      const item = document.createElement('li');
      const link = document.createElement('a');
      const details = document.createElement('small');
      const summary = document.createElement('p');

      link.href = entry.url;
      link.textContent = entry.title;
      details.textContent = `${entry.kind} · ${entry.source} · ${entry.version}`;
      summary.textContent = entry.summary || '';
      item.append(link, details, summary);
      return item;
    });

    results.replaceChildren(...items);
    status.textContent = `${matches.length} result${matches.length === 1 ? '' : 's'}`;
  };

  input.addEventListener('input', render);
  kind.addEventListener('change', render);
  fetch('/search.json')
    .then((response) => response.json())
    .then((data) => { entries = data; render(); })
    .catch(() => { status.textContent = 'Search index could not be loaded.'; });
})();
