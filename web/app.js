fetch('/healthz').then(response => response.json()).then(value => { document.querySelector('#status').textContent = JSON.stringify(value, null, 2) })
