// Theme toggle functionality
document.addEventListener('DOMContentLoaded', () => {
  const themeToggle = document.getElementById('theme-toggle');
  const themeIcon = document.getElementById('theme-icon');
  
  // Check for saved theme preference or use the system preference
  const savedTheme = localStorage.getItem('theme');
  const systemPrefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches;
  
  // Apply the saved theme or system preference
  if (savedTheme === 'dark' || (!savedTheme && systemPrefersDark)) {
    document.documentElement.classList.add('dark');
    themeIcon.classList.replace('fa-moon', 'fa-sun');
  } else {
    document.documentElement.classList.remove('dark');
    themeIcon.classList.replace('fa-sun', 'fa-moon');
  }
  
  // Toggle theme function
  function toggleTheme() {
    if (document.documentElement.classList.contains('dark')) {
      document.documentElement.classList.remove('dark');
      localStorage.setItem('theme', 'light');
      themeIcon.classList.replace('fa-sun', 'fa-moon');
    } else {
      document.documentElement.classList.add('dark');
      localStorage.setItem('theme', 'dark');
      themeIcon.classList.replace('fa-moon', 'fa-sun');
    }
  }
  
  // Add click event to theme toggle button
  themeToggle.addEventListener('click', toggleTheme);
});