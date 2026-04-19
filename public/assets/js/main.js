// Wedding Invitation - Main JavaScript

document.addEventListener('DOMContentLoaded', function() {
  // Initialize countdown
  initCountdown();
  
  // Initialize reveal animations
  initRevealAnimations();
  
  // Initialize music control
  initMusicControl();
});

// Countdown Timer
function initCountdown() {
  const countdown = document.getElementById('countdown');
  if (!countdown) return;
  
  const targetDate = new Date(countdown.dataset.date).getTime();
  
  function updateCountdown() {
    const now = new Date().getTime();
    const distance = targetDate - now;
    
    if (distance < 0) {
      document.getElementById('countdown-days').textContent = '0';
      document.getElementById('countdown-hours').textContent = '0';
      document.getElementById('countdown-minutes').textContent = '0';
      document.getElementById('countdown-seconds').textContent = '0';
      return;
    }
    
    const days = Math.floor(distance / (1000 * 60 * 60 * 24));
    const hours = Math.floor((distance % (1000 * 60 * 60 * 24)) / (1000 * 60 * 60));
    const minutes = Math.floor((distance % (1000 * 60 * 60)) / (1000 * 60));
    const seconds = Math.floor((distance % (1000 * 60)) / 1000);
    
    document.getElementById('countdown-days').textContent = days;
    document.getElementById('countdown-hours').textContent = hours;
    document.getElementById('countdown-minutes').textContent = minutes;
    document.getElementById('countdown-seconds').textContent = seconds;
  }
  
  updateCountdown();
  setInterval(updateCountdown, 1000);
}

// Reveal Animations on Scroll
function initRevealAnimations() {
  const reveals = document.querySelectorAll('.reveal');
  
  function reveal() {
    reveals.forEach(element => {
      const windowHeight = window.innerHeight;
      const elementTop = element.getBoundingClientRect().top;
      const elementVisible = 150;
      
      if (elementTop < windowHeight - elementVisible) {
        element.classList.add('active');
      }
    });
  }
  
  window.addEventListener('scroll', reveal);
  reveal(); // Initial check
}

// Music Control
function initMusicControl() {
  const musicControl = document.getElementById('music-control');
  const audio = document.getElementById('background-music');
  const playIcon = document.querySelector('.play-icon');
  const pauseIcon = document.querySelector('.pause-icon');
  
  if (!musicControl || !audio) return;
  
  let isPlaying = false;
  
  musicControl.addEventListener('click', function() {
    if (isPlaying) {
      audio.pause();
      playIcon.style.display = 'block';
      pauseIcon.style.display = 'none';
    } else {
      audio.play();
      playIcon.style.display = 'none';
      pauseIcon.style.display = 'block';
    }
    isPlaying = !isPlaying;
  });
  
  // Auto-play on first user interaction
  document.addEventListener('click', function autoPlay() {
    if (!isPlaying && audio.src) {
      audio.play().then(() => {
        isPlaying = true;
        playIcon.style.display = 'none';
        pauseIcon.style.display = 'block';
      }).catch(() => {});
    }
    document.removeEventListener('click', autoPlay);
  }, { once: true });
}

// Copy to Clipboard
function copyToClipboard(text, button) {
  navigator.clipboard.writeText(text).then(() => {
    const originalText = button.textContent;
    button.textContent = 'Tersalin!';
    button.style.background = '#28a745';
    
    setTimeout(() => {
      button.textContent = originalText;
      button.style.background = '';
    }, 2000);
  });
}
