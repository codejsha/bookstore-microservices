document.documentElement.style.visibility = "hidden";

function reveal() {
  document.documentElement.style.visibility = "";
}

document.fonts.load('16px "Geist Variable"').then(reveal).catch(reveal);
setTimeout(reveal, 500);
