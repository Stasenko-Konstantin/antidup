function main(): void {
  document.title = "antidup"

  const greeting = document.createElement('h1');
  greeting.textContent = `Hello, World!`;
  document.body.appendChild(greeting);

  const paragraph = document.createElement('p');
  paragraph.textContent = 'This is a simple TypeScript application.';
  document.body.appendChild(paragraph);
}

main();
