const LIMIT = 50
let offset = 0
let loading = false
let hasMore = true

const app = decodeURIComponent(window.location.pathname.replace('/gallery/', ''))
const main = document.getElementById('main')
const lightbox = document.getElementById('lightbox')
const lightboxImg = document.getElementById('lightbox-img')

document.getElementById('title').textContent = app
document.title = `${app} — Bildgalleri`

const sentinel = document.createElement('div')
main.after(sentinel)

const observer = new IntersectionObserver(async entries => {
  if (entries[0].isIntersecting && !loading && hasMore) {
    await loadMore()
  }
}, { rootMargin: '300px' })

observer.observe(sentinel)

async function loadMore() {
  loading = true
  try {
    const res = await fetch(`/api/images?app=${encodeURIComponent(app)}&offset=${offset}&limit=${LIMIT}`)
    if (!res.ok) throw new Error(`HTTP ${res.status}`)
    const data = await res.json()

    if (offset === 0 && data.images.length === 0) {
      main.innerHTML = '<p class="status">Inga bilder uppladdade ännu.</p>'
      observer.disconnect()
      return
    }

    if (offset === 0) {
      main.innerHTML = ''
      const grid = document.createElement('div')
      grid.className = 'grid'
      grid.id = 'grid'
      main.appendChild(grid)
    }

    const grid = document.getElementById('grid')
    for (const { url } of data.images) {
      const img = document.createElement('img')
      img.src = url
      img.loading = 'lazy'
      img.dataset.url = url
      img.alt = ''
      grid.appendChild(img)
    }

    offset += data.images.length
    hasMore = data.hasMore
    if (!hasMore) observer.disconnect()
  } catch (err) {
    if (offset === 0) {
      main.innerHTML = `<p class="status">Kunde inte ladda bilder: ${err.message}</p>`
    }
  } finally {
    loading = false
  }
}

main.addEventListener('click', e => {
  const url = e.target.dataset.url
  if (url) openLightbox(url)
})

lightbox.addEventListener('click', e => {
  if (e.target === lightbox) closeLightbox()
})

document.getElementById('lightbox-close').addEventListener('click', closeLightbox)

document.addEventListener('keydown', e => {
  if (e.key === 'Escape') closeLightbox()
})

function openLightbox(url) {
  lightboxImg.src = url
  lightbox.classList.add('open')
}

function closeLightbox() {
  lightbox.classList.remove('open')
  lightboxImg.src = ''
}
