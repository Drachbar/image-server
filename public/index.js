const main = document.getElementById('main')

async function load() {
  try {
    const res = await fetch('/api/apps')
    if (!res.ok) throw new Error(`HTTP ${res.status}`)
    const apps = await res.json()

    if (!apps || apps.length === 0) {
      main.innerHTML = '<p class="status">Inga bilder uppladdade ännu.</p>'
      return
    }

    main.innerHTML = ''
    const grid = document.createElement('div')
    grid.className = 'apps-grid'

    for (const app of apps) {
      const card = document.createElement('a')
      card.className = 'app-card'
      card.href = `/gallery/${encodeURIComponent(app.name)}`
      card.innerHTML = `
        <img src="${app.thumbnail}" loading="lazy" alt="${app.name}">
        <div class="app-card-info">
          <h2>${app.name}</h2>
          <span>${app.count} ${app.count === 1 ? 'bild' : 'bilder'}</span>
        </div>
      `
      grid.appendChild(card)
    }

    main.appendChild(grid)
  } catch (err) {
    main.innerHTML = `<p class="status">Kunde inte ladda appar: ${err.message}</p>`
  }
}

load()
