function loadZip(event: Event)  {
  const target = event.target as HTMLInputElement
  if (!target.files || target.files.length != 1) {
    alert("need one file!")
    return
  }
  const file = target.files[0]
  if (file.type != "application/zip") {
    alert("need only .zip!")
  }
}

function main(): void {
  document.title = "antidup"

  const title = document.createElement('h1');
  title.textContent = `antidup`;
  document.body.appendChild(title);

  const paragraph = document.createElement('p');
  paragraph.textContent = 'load .zip archive of your pictures and I say what duplicates it have.';
  document.body.appendChild(paragraph);

  const fileInput = document.createElement("input")
  fileInput.type = "file"
  fileInput.style.display = "none"
  fileInput.addEventListener("change", loadZip)
  document.body.appendChild(fileInput)

  const loadButton = document.createElement("button")
  loadButton.textContent = "load .zip"
  loadButton.addEventListener("click", () => {
    fileInput.click()
  })
  document.body.appendChild(loadButton)

  document.body.appendChild(document.createElement("hr"))

  const author = document.createElement("a")
  author.href = "https://github.com/Stasenko-Konstantin"
  author.textContent = "author github page"
  author.target = "_blank"
  document.body.appendChild(author)
}

main();
