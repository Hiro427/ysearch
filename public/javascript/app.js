document.body.addEventListener("htmx:afterSwap", () => {
  if (window.Alpine) window.Alpine.initTree(document.body);
});

document.addEventListener("alpine:init", () => {
  Alpine.store("theme", {
    dark: localStorage.getItem("theme") === "dark",

    toggle() {
      this.dark = !this.dark;
      localStorage.setItem("theme", this.dark ? "dark" : "light");
    },
  });
  Alpine.data("raceList", () => ({
    races: [
      "Aqua",
      "Beast",
      "Beast-Warrior",
      "Creator God",
      "Cyberse",
      "Dinosaur",
      "Divine-Beast",
      "Dragon",
      "Fairy",
      "Fiend",
      "Fish",
      "Insect",
      "Illusion",
      "Machine",
      "Plant",
      "Psychic",
      "Pyro",
      "Reptile",
      "Rock",
      "Sea Serpent",
      "Spellcaster",
      "Thunder",
      "Warrior",
      "Winged Beast",
      "Wyrm",
      "Zombie",
    ],
  }));
});
