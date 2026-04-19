package seed_themes

import (
	"fmt"
	"invitation-api/internal/database"
	"invitation-api/internal/domain/theme"

	"github.com/google/uuid"
)

func CreateInitialCategories() {
	db := database.GetDB()

	categories := []theme.ThemeCategory{
		{
			ID:          uuid.New(),
			Name:        "One Page Scroll",
			Slug:        "one-page",
			Description: "Single page scrolling invitation, perfect for simple and elegant designs.",
			DefaultHTML: `<!DOCTYPE html>
<html lang="id" class="scroll-smooth">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>The Wedding of {{.GroomName}} & {{.BrideName}}</title>
    <!-- Tailwind CSS -->
    <script src="https://cdn.tailwindcss.com"></script>
    <script>
        tailwind.config = {
            theme: {
                extend: {
                    colors: {
                        primary: '#C5A880',
                        secondary: '#5A4633',
                        dark: '#1a1a1a',
                    },
                    fontFamily: {
                        serif: ['Playfair Display', 'serif'],
                        sans: ['Lato', 'sans-serif'],
                        cursive: ['Great Vibes', 'cursive'],
                    },
                }
            }
        }
    </script>
    <!-- Google Fonts -->
    <link href="https://fonts.googleapis.com/css2?family=Playfair+Display:ital,wght@0,400;0,600;1,400&family=Lato:wght@300;400&family=Great+Vibes&display=swap" rel="stylesheet">
    <style>
        .animate-fade-in-up { animation: fadeInUp 1s ease-out forwards; opacity: 0; }
        .delay-100 { animation-delay: 0.1s; }
        .delay-200 { animation-delay: 0.2s; }
        .delay-300 { animation-delay: 0.3s; }
        @keyframes fadeInUp { from { opacity: 0; transform: translateY(20px); } to { opacity: 1; transform: translateY(0); } }
        
        {{if .BackgroundImage}}
        body {
            background-image: url('{{.BackgroundImage}}');
            background-size: cover;
            background-attachment: fixed;
            background-position: center;
        }
        {{else}}
        body {
            background-color: #fcfcfc;
        }
        {{end}}
        
        .glass-card {
            background: rgba(255, 255, 255, 0.9);
            backdrop-filter: blur(10px);
            border-radius: 1rem;
            box-shadow: 0 4px 30px rgba(0, 0, 0, 0.1);
            border: 1px solid rgba(255, 255, 255, 0.3);
        }
    </style>
</head>
<body class="font-sans antialiased text-gray-800">

    <!-- Cover Section -->
    <section id="cover" class="h-screen flex flex-col items-center justify-center text-center bg-cover bg-center relative" style="background-image: url('{{if .CoverImage}}{{.CoverImage}}{{else}}{{.CouplePhoto}}{{end}}');">
        <div class="absolute inset-0 bg-black/50"></div>
        <div class="relative z-10 text-white p-6 max-w-2xl mx-auto">
            <h3 class="text-xl tracking-[0.2em] uppercase mb-4 animate-fade-in-up">The Wedding Of</h3>
            <h1 class="font-cursive text-6xl md:text-8xl mb-6 animate-fade-in-up delay-100">{{.GroomNickname}} & {{.BrideNickname}}</h1>
            <p class="text-xl font-serif italic mb-8 animate-fade-in-up delay-200">{{.AkadDateFormatted}}</p>
            
            <div class="glass-card p-6 text-gray-800 inline-block animate-fade-in-up delay-300 transform hover:scale-105 transition duration-300">
                <p class="mb-2 text-sm uppercase tracking-wide">Kepada Yth:</p>
                <p class="text-lg font-bold mb-4">Tamu Undangan</p>
                <button onclick="openInvitation()" class="px-8 py-3 bg-secondary text-white rounded-full hover:bg-opacity-90 transition shadow-lg flex items-center gap-2 mx-auto">
                    <svg xmlns="http://www.w3.org/2000/svg" class="h-5 w-5" viewBox="0 0 20 20" fill="currentColor">
                        <path d="M10 2a6 6 0 00-6 6v3.586l-.707.707A1 1 0 004 14h12a1 1 0 00.707-1.707L16 11.586V8a6 6 0 00-6-6zM10 18a3 3 0 01-3-3h6a3 3 0 01-3 3z" />
                    </svg>
                    Buka Undangan
                </button>
            </div>
        </div>
    </section>

    <!-- Content (Hidden initially, shown after open) -->
    <div id="content" class="hidden">
    
        <!-- Music Player -->
        {{if .MusicURL}}
        <div id="music-control" class="fixed bottom-6 right-6 z-50 animate-bounce">
            <button onclick="toggleMusic()" class="p-3 bg-white/90 backdrop-blur rounded-full shadow-xl border border-secondary text-secondary hover:scale-110 transition">
                <svg id="icon-play" class="w-6 h-6 hidden" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M14.752 11.168l-3.197-2.132A1 1 0 0010 9.87v4.263a1 1 0 001.555.832l3.197-2.132a1 1 0 000-1.664z"></path><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path></svg>
                <svg id="icon-pause" class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 9v6m4-6v6m7-3a9 9 0 11-18 0 9 9 0 0118 0z"></path></svg>
            </button>
        </div>
        <audio id="bgm" loop>
            <source src="{{.MusicURL}}" type="audio/mp3">
        </audio>
        {{end}}

        <!-- Couple Section -->
        <section class="py-24 px-6 text-center">
            <h2 class="font-cursive text-5xl text-secondary mb-16">Mempelai</h2>
            
            <div class="flex flex-col md:flex-row justify-center gap-16 items-center max-w-4xl mx-auto">
                <!-- Groom -->
                <div class="flex-1">
                    <div class="relative w-48 h-48 mx-auto mb-6 rounded-full overflow-hidden border-4 border-secondary shadow-xl">
                         <img src="{{if .GroomPhoto}}{{.GroomPhoto}}{{else}}https://via.placeholder.com/300{{end}}" class="w-full h-full object-cover">
                    </div>
                    <h3 class="font-serif text-3xl mb-2 font-bold text-gray-800">{{.GroomName}}</h3>
                    <p class="text-gray-600 mb-4 text-xl">{{.GroomNickname}}</p>
                    <p class="text-sm text-gray-500">Putra dari Bpk. {{.GroomFather}} & Ibu {{.GroomMother}}</p>
                </div>
                
                <div class="font-cursive text-6xl text-primary">&</div>
                
                <!-- Bride -->
                <div class="flex-1">
                    <div class="relative w-48 h-48 mx-auto mb-6 rounded-full overflow-hidden border-4 border-secondary shadow-xl">
                        <img src="{{if .BridePhoto}}{{.BridePhoto}}{{else}}https://via.placeholder.com/300{{end}}" class="w-full h-full object-cover">
                    </div>
                    <h3 class="font-serif text-3xl mb-2 font-bold text-gray-800">{{.BrideName}}</h3>
                    <p class="text-gray-600 mb-4 text-xl">{{.BrideNickname}}</p>
                    <p class="text-sm text-gray-500">Putri dari Bpk. {{.BrideFather}} & Ibu {{.BrideMother}}</p>
                </div>
            </div>
        </section>

        <!-- Events Section -->
        <section class="py-20 bg-secondary/5 px-6 text-center">
            <h2 class="font-cursive text-5xl text-secondary mb-12">Acara</h2>
            <div class="grid grid-cols-1 md:grid-cols-2 gap-8 max-w-4xl mx-auto">
                 <div class="glass-card p-8">
                     <h3 class="font-serif text-3xl mb-4">Akad Nikah</h3>
                     <p class="text-xl font-bold mb-2">{{.AkadDateFormatted}}</p>
                     <p class="mb-6 text-gray-600">{{.AkadTime}}</p>
                     <p class="text-gray-600 mb-6 font-medium">{{.AkadLocation}}</p>
                     <p class="text-gray-500 text-sm mb-6">{{.AkadAddress}}</p>
                     <a href="{{.AkadMapsURL}}" target="_blank" class="inline-block px-6 py-2 bg-secondary text-white rounded hover:bg-secondary/90 transition">Lihat Lokasi</a>
                 </div>
                 
                 <div class="glass-card p-8">
                     <h3 class="font-serif text-3xl mb-4">Resepsi</h3>
                     <p class="text-xl font-bold mb-2">{{if .ReceptionDate}}{{.ReceptionDate}}{{else}}{{.AkadDateFormatted}}{{end}}</p>
                     <p class="mb-6 text-gray-600">{{.ReceptionTime}}</p>
                     <p class="text-gray-600 mb-6 font-medium">{{.ReceptionLocation}}</p>
                     <p class="text-gray-500 text-sm mb-6">{{.ReceptionAddress}}</p>
                     <a href="{{.ReceptionMapsURL}}" target="_blank" class="inline-block px-6 py-2 bg-secondary text-white rounded hover:bg-secondary/90 transition">Lihat Lokasi</a>
                 </div>
            </div>
        </section>
        
        <!-- Gallery Section -->
        <section class="py-20 px-6">
            <h2 class="font-cursive text-5xl text-secondary mb-12 text-center">Galeri Foto</h2>
            <div class="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4 max-w-6xl mx-auto">
                {{range .Gallery}}
                <div class="aspect-square overflow-hidden rounded-xl shadow-md hover:shadow-xl transition duration-300 transform hover:scale-105 cursor-pointer group">
                    <img src="{{.URL}}" class="w-full h-full object-cover group-hover:opacity-90 transition">
                </div>
                {{end}}
            </div>
        </section>

        <!-- RSVP & Wishes -->
        <section class="py-20 px-6 bg-stone-100">
            <div class="max-w-xl mx-auto">
                <h2 class="font-cursive text-5xl text-secondary mb-12 text-center">Ucapan & Doa</h2>
                 <!-- HTMX Messages -->
                 <div hx-get="/invitations/{{.Subdomain}}/messages" hx-trigger="load"></div>
                 
                 <!-- HTMX RSVP Form -->
                 <div class="mt-12">
                      <div hx-get="/invitations/{{.Subdomain}}/rsvp" hx-trigger="load"></div>
                 </div>
            </div>
        </section>
        
        <!-- Footer -->
        <footer class="py-10 text-center text-gray-500 text-sm bg-gray-900 text-white">
            <p>&copy; 2025 {{.GroomName}} & {{.BrideName}}</p>
        </footer>
    </div>

    <script>
        function openInvitation() {
            document.getElementById('content').classList.remove('hidden');
            // document.getElementById('cover').scrollIntoView({ behavior: 'smooth' });
            
            // Play Music
            const audio = document.getElementById('bgm');
            if(audio) {
                audio.play().catch(e => console.log("Audio autoplay failed", e));
            }
            
            // Scroll to content
            setTimeout(() => {
                document.getElementById('content').scrollIntoView({ behavior: 'smooth' });
            }, 300);
        }

        function toggleMusic() {
            const audio = document.getElementById('bgm');
            const iconPlay = document.getElementById('icon-play');
            const iconPause = document.getElementById('icon-pause');
            
            if (audio.paused) {
                audio.play();
                iconPlay.classList.add('hidden');
                iconPause.classList.remove('hidden');
            } else {
                audio.pause();
                iconPlay.classList.remove('hidden');
                iconPause.classList.add('hidden');
            }
        }
    </script>
</body>
</html>`,
			DefaultCSS: `/* Custom Global CSS */
body {
    -webkit-font-smoothing: antialiased;
}
.font-serif {
    font-family: 'Playfair Display', serif;
}`,
		},
		{
			ID:          uuid.New(),
			Name:        "Modern Side Menu",
			Slug:        "modern-side",
			Description: "Modern layout with a fixed sidebar navigation.",
			DefaultHTML: `<!DOCTYPE html>
<html lang="id" class="scroll-smooth">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>{{.GroomName}} & {{.BrideName}}</title>
    <script src="https://cdn.tailwindcss.com"></script>
    <link href="https://fonts.googleapis.com/css2?family=Outfit:wght@300;500;700&family=Great+Vibes&display=swap" rel="stylesheet">
    <script>
        tailwind.config = {
            theme: {
                extend: {
                    colors: {
                        primary: '#4F46E5',
                        secondary: '#818CF8',
                    },
                    fontFamily: {
                        sans: ['Outfit', 'sans-serif'],
                        cursive: ['Great Vibes', 'cursive'],
                    }
                }
            }
        }
    </script>
    <style>
        .sidebar-active { transform: translateX(0); }
        {{if .BackgroundImage}}
        .main-bg {
            background-image: url('{{.BackgroundImage}}');
            background-size: cover;
            background-attachment: fixed;
            background-position: center;
        }
        {{end}}
    </style>
</head>
<body class="bg-gray-50 text-gray-800 lg:pl-80">

    <!-- Music Player -->
    {{if .MusicURL}}
    <div id="music-control" class="fixed bottom-6 right-6 z-50">
        <button onclick="toggleMusic()" class="p-3 bg-white rounded-full shadow-xl border border-primary text-primary hover:scale-110 transition">
            <svg id="icon-play" class="w-6 h-6 hidden" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M14.752 11.168l-3.197-2.132A1 1 0 0010 9.87v4.263a1 1 0 001.555.832l3.197-2.132a1 1 0 000-1.664z"></path><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 12a9 9 0 11-18 0 9 9 0 0118 0z"></path></svg>
            <svg id="icon-pause" class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24"><path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 9v6m4-6v6m7-3a9 9 0 11-18 0 9 9 0 0118 0z"></path></svg>
        </button>
    </div>
    <audio id="bgm" loop autoplay>
        <source src="{{.MusicURL}}" type="audio/mp3">
    </audio>
    {{end}}

    <!-- Sidebar Navigation -->
    <aside class="fixed inset-y-0 left-0 w-80 bg-white border-r border-gray-200 p-8 hidden lg:flex flex-col justify-between z-50">
        <div>
            <h1 class="text-3xl font-bold tracking-tighter mb-2">{{.GroomNickname}}</h1>
            <h1 class="text-3xl font-bold tracking-tighter text-primary">& {{.BrideNickname}}</h1>
            <p class="mt-4 text-gray-500">{{.AkadDateFormatted}}</p>
        </div>
        
        <nav class="space-y-4 text-lg font-medium">
            <a href="#home" class="block hover:text-primary transition">Beranda</a>
            <a href="#couple" class="block hover:text-primary transition">Pasangan</a>
            <a href="#event" class="block hover:text-primary transition">Acara</a>
            <a href="#gallery" class="block hover:text-primary transition">Galeri</a>
            <a href="#wishes" class="block hover:text-primary transition">Ucapan</a>
        </nav>

        <div class="text-sm text-gray-400">
            &copy; 2025 Invitation Online
        </div>
    </aside>

    <!-- Mobile Header -->
    <header class="lg:hidden fixed top-0 w-full bg-white/80 backdrop-blur border-b z-40 px-6 py-4 flex justify-between items-center">
        <span class="font-bold text-xl">{{.GroomNickname}} & {{.BrideNickname}}</span>
        <button onclick="toggleMobileMenu()" class="text-2xl">☰</button>
    </header>

    <!-- Content -->
    <main class="pt-16 lg:pt-0 main-bg">
        <!-- Home Section -->
        <section id="home" class="min-h-screen flex items-center justify-center bg-gradient-to-br from-primary/5 to-secondary/10 px-6">
            <div class="text-center max-w-lg">
                <span class="inline-block px-3 py-1 bg-primary/10 text-primary rounded-full text-sm font-medium mb-6">Save The Date</span>
                <h1 class="text-5xl md:text-7xl font-bold mb-6 leading-tight">We Are<br>Getting Married</h1>
                <p class="text-xl text-gray-600 mb-8">Kami mengundang Anda untuk merayakan cinta kami.</p>
                <img src="{{if .CoverImage}}{{.CoverImage}}{{else}}{{.CouplePhoto}}{{end}}" class="w-full h-64 object-cover rounded-2xl shadow-2xl mb-8 hover:scale-105 transition duration-500">
            </div>
        </section>

        <!-- Couple Section -->
        <section id="couple" class="py-24 px-6 bg-white">
            <h2 class="font-cursive text-5xl text-center text-primary mb-16">Pasangan</h2>
            <div class="flex flex-col md:flex-row justify-center gap-16 items-center max-w-4xl mx-auto">
                <div class="flex-1 text-center">
                    <div class="w-48 h-48 mx-auto mb-6 rounded-full overflow-hidden border-4 border-primary shadow-xl">
                         <img src="{{if .GroomPhoto}}{{.GroomPhoto}}{{else}}https://via.placeholder.com/300{{end}}" class="w-full h-full object-cover">
                    </div>
                    <h3 class="text-2xl font-bold mb-2">{{.GroomName}}</h3>
                    <p class="text-gray-500">Putra dari Bpk. {{.GroomFather}} & Ibu {{.GroomMother}}</p>
                </div>
                
                <div class="font-cursive text-6xl text-primary">&</div>
                
                <div class="flex-1 text-center">
                    <div class="w-48 h-48 mx-auto mb-6 rounded-full overflow-hidden border-4 border-primary shadow-xl">
                        <img src="{{if .BridePhoto}}{{.BridePhoto}}{{else}}https://via.placeholder.com/300{{end}}" class="w-full h-full object-cover">
                    </div>
                    <h3 class="text-2xl font-bold mb-2">{{.BrideName}}</h3>
                    <p class="text-gray-500">Putri dari Bpk. {{.BrideFather}} & Ibu {{.BrideMother}}</p>
                </div>
            </div>
        </section>
        
        <!-- Event Section -->
        <section id="event" class="py-24 px-6 bg-gray-50">
            <h2 class="font-cursive text-5xl text-center text-primary mb-16">Acara</h2>
            <div class="grid grid-cols-1 md:grid-cols-2 gap-8 max-w-4xl mx-auto">
                <div class="bg-white p-8 rounded-2xl shadow-lg hover:shadow-xl transition">
                    <h3 class="text-2xl font-bold mb-4 text-primary">Akad Nikah</h3>
                    <p class="text-lg font-medium mb-2">{{.AkadDateFormatted}}</p>
                    <p class="text-gray-600 mb-4">{{.AkadTime}}</p>
                    <p class="text-gray-500 mb-6">{{.AkadAddress}}</p>
                    <a href="{{.AkadMapsURL}}" target="_blank" class="inline-block px-6 py-2 bg-primary text-white rounded-lg hover:bg-primary/90 transition">Lihat Lokasi</a>
                </div>
                <div class="bg-white p-8 rounded-2xl shadow-lg hover:shadow-xl transition">
                    <h3 class="text-2xl font-bold mb-4 text-primary">Resepsi</h3>
                    <p class="text-lg font-medium mb-2">{{if .ReceptionDate}}{{.ReceptionDate}}{{else}}{{.AkadDateFormatted}}{{end}}</p>
                    <p class="text-gray-600 mb-4">{{.ReceptionTime}}</p>
                    <p class="text-gray-500 mb-6">{{.ReceptionAddress}}</p>
                    <a href="{{.ReceptionMapsURL}}" target="_blank" class="inline-block px-6 py-2 bg-primary text-white rounded-lg hover:bg-primary/90 transition">Lihat Lokasi</a>
                </div>
            </div>
        </section>

        <!-- Gallery Section -->
        <section id="gallery" class="py-24 px-6 bg-white">
            <h2 class="font-cursive text-5xl text-center text-primary mb-16">Galeri</h2>
            <div class="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4 max-w-6xl mx-auto">
                {{range .Gallery}}
                <div class="aspect-square overflow-hidden rounded-xl shadow-md hover:shadow-xl transition">
                    <img src="{{.URL}}" class="w-full h-full object-cover hover:scale-110 transition duration-500">
                </div>
                {{end}}
            </div>
        </section>

        <!-- Wishes Section -->
        <section id="wishes" class="py-24 px-6 bg-gray-50">
            <h2 class="font-cursive text-5xl text-center text-primary mb-16">Ucapan & Doa</h2>
            <div class="max-w-xl mx-auto">
                <div hx-get="/invitations/{{.Subdomain}}/messages" hx-trigger="load"></div>
                <div class="mt-8" hx-get="/invitations/{{.Subdomain}}/rsvp" hx-trigger="load"></div>
            </div>
        </section>

        <!-- Footer -->
        <footer class="py-10 text-center bg-gray-900 text-white">
            <p>&copy; 2025 {{.GroomName}} & {{.BrideName}}</p>
        </footer>
    </main>
    
    <script>
        function toggleMusic() {
            const audio = document.getElementById('bgm');
            const iconPlay = document.getElementById('icon-play');
            const iconPause = document.getElementById('icon-pause');
            
            if (audio.paused) {
                audio.play();
                iconPlay.classList.add('hidden');
                iconPause.classList.remove('hidden');
            } else {
                audio.pause();
                iconPlay.classList.remove('hidden');
                iconPause.classList.add('hidden');
            }
        }
        
        function toggleMobileMenu() {
            // Simple mobile menu toggle
            alert('Mobile menu - implement based on your design');
        }
    </script>
</body>
</html>`,
			DefaultCSS: `/* Modern CSS */
html {
    scroll-behavior: smooth;
}`,
		},
		{
			ID:          uuid.New(),
			Name:        "Luxury Premium",
			Slug:        "luxury",
			Description: "Luxury premium invitation with slideshow, countdown, love story, bank gift, and premium animations.",
			DefaultHTML: `<!DOCTYPE html>
<html lang="id" class="scroll-smooth">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0, maximum-scale=1.0, user-scalable=no">
    <title>{{.GroomName}} & {{.BrideName}} - Wedding Invitation</title>
    
    <!-- Tailwind CSS -->
    <script src="https://cdn.tailwindcss.com"></script>
    <script>
        tailwind.config = {
            theme: {
                extend: {
                    colors: {
                        primary: '#d2d0cc',
                        secondary: '#2A2A2A',
                        gold: '#C5A880',
                        cream: '#FAF8F5',
                    },
                    fontFamily: {
                        serif: ['Playfair Display', 'serif'],
                        sans: ['Poppins', 'sans-serif'],
                        cursive: ['Great Vibes', 'cursive'],
                        display: ['Italiana', 'serif'],
                    }
                }
            }
        }
    </script>
    
    <!-- Google Fonts -->
    <link href="https://fonts.googleapis.com/css2?family=Playfair+Display:ital,wght@0,400;0,600;1,400&family=Poppins:wght@300;400;500;600&family=Great+Vibes&family=Italiana&display=swap" rel="stylesheet">
    
    <style>
        /* Animations */
        @keyframes fadeInUp {
            from { opacity: 0; transform: translateY(30px); }
            to { opacity: 1; transform: translateY(0); }
        }
        @keyframes fadeInDown {
            from { opacity: 0; transform: translateY(-30px); }
            to { opacity: 1; transform: translateY(0); }
        }
        @keyframes zoomIn {
            from { opacity: 0; transform: scale(0.8); }
            to { opacity: 1; transform: scale(1); }
        }
        @keyframes bounce {
            0%, 100% { transform: translateY(0); }
            50% { transform: translateY(-10px); }
        }
        @keyframes pulse {
            0%, 100% { opacity: 1; }
            50% { opacity: 0.5; }
        }
        @keyframes kenburns {
            0% { transform: scale(1); }
            100% { transform: scale(1.1); }
        }
        
        .animate-fade-in-up { animation: fadeInUp 1s ease-out forwards; }
        .animate-fade-in-down { animation: fadeInDown 1s ease-out forwards; }
        .animate-zoom-in { animation: zoomIn 0.8s ease-out forwards; }
        .animate-bounce-slow { animation: bounce 2s ease-in-out infinite; }
        .animate-pulse-slow { animation: pulse 3s ease-in-out infinite; }
        .animate-kenburns { animation: kenburns 20s ease-in-out infinite alternate; }
        
        .delay-100 { animation-delay: 0.1s; }
        .delay-200 { animation-delay: 0.2s; }
        .delay-300 { animation-delay: 0.3s; }
        .delay-400 { animation-delay: 0.4s; }
        
        /* Glass Effect */
        .glass {
            background: rgba(255, 255, 255, 0.85);
            backdrop-filter: blur(10px);
            -webkit-backdrop-filter: blur(10px);
        }
        
        /* Cover Background */
        .cover-bg {
            background-image: url('{{if .CoverImage}}{{.CoverImage}}{{else}}{{.CouplePhoto}}{{end}}');
            background-size: cover;
            background-position: center;
        }
        
        /* Timeline */
        .timeline-line {
            position: absolute;
            left: 50%;
            top: 0;
            bottom: 0;
            width: 2px;
            background: linear-gradient(to bottom, #C5A880, #d2d0cc);
            transform: translateX(-50%);
        }
        
        /* Bank Card Gradient */
        .bank-bca { background: linear-gradient(135deg, #0066AE 0%, #004D8C 100%); }
        .bank-bri { background: linear-gradient(135deg, #00529C 0%, #003D73 100%); }
        .bank-dana { background: linear-gradient(135deg, #118EEA 0%, #0078D4 100%); }
        
        /* Hide scrollbar */
        .no-scrollbar::-webkit-scrollbar { display: none; }
        .no-scrollbar { -ms-overflow-style: none; scrollbar-width: none; }
        
        /* Section backgrounds */
        {{if .BackgroundImage}}
        .main-bg {
            background-image: url('{{.BackgroundImage}}');
            background-size: cover;
            background-attachment: fixed;
            background-position: center;
        }
        {{end}}
    </style>
</head>
<body class="font-sans antialiased text-secondary bg-cream overflow-x-hidden">

    <!-- Opening Cover (Before Open) -->
    <section id="opening" class="fixed inset-0 z-50 cover-bg flex flex-col items-center justify-center text-center">
        <div class="absolute inset-0 bg-black/50"></div>
        <div class="absolute inset-0 animate-kenburns cover-bg opacity-30"></div>
        
        <div class="relative z-10 text-white p-6 max-w-lg mx-auto">
            <p class="text-lg tracking-[0.3em] uppercase mb-4 animate-fade-in-down opacity-0" style="animation-delay: 0.2s;">The Wedding Of</p>
            <h1 class="font-cursive text-6xl md:text-8xl mb-6 animate-zoom-in opacity-0" style="animation-delay: 0.4s;">{{.GroomNickname}} & {{.BrideNickname}}</h1>
            <p class="text-xl font-serif italic mb-12 animate-fade-in-up opacity-0" style="animation-delay: 0.6s;">{{.AkadDateFormatted}}</p>
            
            <div class="glass rounded-2xl p-6 text-secondary animate-fade-in-up opacity-0" style="animation-delay: 0.8s;">
                <p class="text-sm uppercase tracking-widest text-gray-500 mb-2">Kepada Yth.</p>
                <p class="text-xl font-semibold mb-6">Tamu Undangan</p>
                <button onclick="openInvitation()" class="px-8 py-3 bg-gold text-white rounded-full hover:bg-gold/90 transition-all duration-300 shadow-lg hover:shadow-xl flex items-center gap-3 mx-auto animate-bounce-slow">
                    <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 19v-8.93a2 2 0 01.89-1.664l7-4.666a2 2 0 012.22 0l7 4.666A2 2 0 0121 10.07V19M3 19a2 2 0 002 2h14a2 2 0 002-2M3 19l6.75-4.5M21 19l-6.75-4.5M3 10l6.75 4.5M21 10l-6.75 4.5m0 0l-1.14.76a2 2 0 01-2.22 0l-1.14-.76"/>
                    </svg>
                    Buka Undangan
                </button>
            </div>
        </div>
    </section>

    <!-- Main Content (Hidden until opened) -->
    <main id="content" class="hidden main-bg">
    
        <!-- Music Player -->
        {{if .MusicURL}}
        <div id="music-control" class="fixed bottom-6 right-6 z-50">
            <button onclick="toggleMusic()" class="w-12 h-12 bg-gold text-white rounded-full shadow-xl flex items-center justify-center hover:scale-110 transition-transform duration-300 animate-pulse-slow">
                <svg id="music-icon" class="w-6 h-6" fill="currentColor" viewBox="0 0 24 24">
                    <path d="M12 3v10.55c-.59-.34-1.27-.55-2-.55-2.21 0-4 1.79-4 4s1.79 4 4 4 4-1.79 4-4V7h4V3h-6z"/>
                </svg>
            </button>
        </div>
        <audio id="bgm" loop>
            <source src="{{.MusicURL}}" type="audio/mp3">
        </audio>
        {{end}}

        <!-- Hero Section with Countdown -->
        <section class="min-h-screen cover-bg relative flex items-center justify-center text-center">
            <div class="absolute inset-0 bg-gradient-to-b from-black/60 via-black/40 to-black/60"></div>
            
            <div class="relative z-10 text-white p-6 max-w-2xl mx-auto">
                <p class="text-lg tracking-[0.2em] uppercase mb-4">The Wedding Of</p>
                <h1 class="font-cursive text-6xl md:text-8xl mb-4">{{.GroomNickname}} & {{.BrideNickname}}</h1>
                <p class="font-serif text-xl italic mb-12">"And among His signs is that He created for you mates from among yourselves."</p>
                
                <!-- Countdown -->
                <div id="countdown" class="grid grid-cols-4 gap-4 mb-8">
                    <div class="glass rounded-xl p-4 text-secondary">
                        <span id="days" class="block text-3xl font-bold text-gold">00</span>
                        <span class="text-sm">Hari</span>
                    </div>
                    <div class="glass rounded-xl p-4 text-secondary">
                        <span id="hours" class="block text-3xl font-bold text-gold">00</span>
                        <span class="text-sm">Jam</span>
                    </div>
                    <div class="glass rounded-xl p-4 text-secondary">
                        <span id="minutes" class="block text-3xl font-bold text-gold">00</span>
                        <span class="text-sm">Menit</span>
                    </div>
                    <div class="glass rounded-xl p-4 text-secondary">
                        <span id="seconds" class="block text-3xl font-bold text-gold">00</span>
                        <span class="text-sm">Detik</span>
                    </div>
                </div>
                
                <a href="https://calendar.google.com/calendar/render?action=TEMPLATE&text=Pernikahan+{{.GroomNickname}}+dan+{{.BrideNickname}}" target="_blank" class="inline-flex items-center gap-2 px-6 py-3 border-2 border-white text-white rounded-full hover:bg-white hover:text-secondary transition-all duration-300">
                    <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z"/>
                    </svg>
                    Save The Date
                </a>
            </div>
        </section>

        <!-- Bride & Groom Section -->
        <section class="py-24 px-6 bg-primary/30">
            <div class="max-w-4xl mx-auto text-center">
                <h2 class="font-cursive text-5xl text-gold mb-4">Bride & Groom</h2>
                <p class="text-gray-600 mb-16 max-w-xl mx-auto">Dengan memohon rahmat dan ridho Allah SWT, kami bermaksud menyelenggarakan pernikahan putra-putri kami</p>
                
                <div class="flex flex-col md:flex-row justify-center items-center gap-12">
                    <!-- Groom -->
                    <div class="flex-1 max-w-xs">
                        <div class="w-48 h-48 mx-auto mb-6 rounded-full overflow-hidden border-4 border-gold shadow-2xl">
                            <img src="{{if .GroomPhoto}}{{.GroomPhoto}}{{else}}https://via.placeholder.com/300{{end}}" class="w-full h-full object-cover" alt="{{.GroomName}}">
                        </div>
                        <h3 class="font-cursive text-4xl text-gold mb-2">{{.GroomName}}</h3>
                        <p class="text-gray-600 mb-4">{{.GroomNickname}}</p>
                        <p class="text-sm text-gray-500">Putra dari<br>Bpk. {{.GroomFather}} & Ibu {{.GroomMother}}</p>
                    </div>
                    
                    <div class="font-cursive text-6xl text-gold">&</div>
                    
                    <!-- Bride -->
                    <div class="flex-1 max-w-xs">
                        <div class="w-48 h-48 mx-auto mb-6 rounded-full overflow-hidden border-4 border-gold shadow-2xl">
                            <img src="{{if .BridePhoto}}{{.BridePhoto}}{{else}}https://via.placeholder.com/300{{end}}" class="w-full h-full object-cover" alt="{{.BrideName}}">
                        </div>
                        <h3 class="font-cursive text-4xl text-gold mb-2">{{.BrideName}}</h3>
                        <p class="text-gray-600 mb-4">{{.BrideNickname}}</p>
                        <p class="text-sm text-gray-500">Putri dari<br>Bpk. {{.BrideFather}} & Ibu {{.BrideMother}}</p>
                    </div>
                </div>
            </div>
        </section>

        <!-- Events Section -->
        <section class="py-24 px-6 cover-bg relative">
            <div class="absolute inset-0 bg-black/70"></div>
            <div class="relative z-10 max-w-4xl mx-auto">
                <h2 class="font-cursive text-5xl text-gold text-center mb-16">Waktu & Tempat</h2>
                
                <div class="grid grid-cols-1 md:grid-cols-2 gap-8">
                    <!-- Akad -->
                    <div class="glass rounded-2xl p-8 text-center">
                        <h3 class="font-display text-3xl text-gold mb-6">Akad Nikah</h3>
                        <div class="text-6xl font-bold text-gold mb-2">{{.AkadDay}}</div>
                        <p class="text-xl font-semibold mb-2">{{.AkadDateFormatted}}</p>
                        <p class="text-gray-600 mb-6">{{.AkadTime}}</p>
                        <div class="border-t border-gray-200 pt-6">
                            <p class="font-semibold mb-2">{{.AkadLocation}}</p>
                            <p class="text-sm text-gray-500 mb-4">{{.AkadAddress}}</p>
                            <a href="{{.AkadMapsURL}}" target="_blank" class="inline-flex items-center gap-2 px-4 py-2 bg-gold text-white rounded-lg hover:bg-gold/90 transition">
                                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17.657 16.657L13.414 20.9a1.998 1.998 0 01-2.827 0l-4.244-4.243a8 8 0 1111.314 0z"/>
                                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 11a3 3 0 11-6 0 3 3 0 016 0z"/>
                                </svg>
                                Lihat Lokasi
                            </a>
                        </div>
                    </div>
                    
                    <!-- Resepsi -->
                    <div class="glass rounded-2xl p-8 text-center">
                        <h3 class="font-display text-3xl text-gold mb-6">Resepsi</h3>
                        <div class="text-6xl font-bold text-gold mb-2">{{.ReceptionDay}}</div>
                        <p class="text-xl font-semibold mb-2">{{if .ReceptionDate}}{{.ReceptionDate}}{{else}}{{.AkadDateFormatted}}{{end}}</p>
                        <p class="text-gray-600 mb-6">{{.ReceptionTime}}</p>
                        <div class="border-t border-gray-200 pt-6">
                            <p class="font-semibold mb-2">{{.ReceptionLocation}}</p>
                            <p class="text-sm text-gray-500 mb-4">{{.ReceptionAddress}}</p>
                            <a href="{{.ReceptionMapsURL}}" target="_blank" class="inline-flex items-center gap-2 px-4 py-2 bg-gold text-white rounded-lg hover:bg-gold/90 transition">
                                <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17.657 16.657L13.414 20.9a1.998 1.998 0 01-2.827 0l-4.244-4.243a8 8 0 1111.314 0z"/>
                                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 11a3 3 0 11-6 0 3 3 0 016 0z"/>
                                </svg>
                                Lihat Lokasi
                            </a>
                        </div>
                    </div>
                </div>
            </div>
        </section>

        <!-- Love Story Section -->
        <section class="py-24 px-6 bg-cream">
            <div class="max-w-3xl mx-auto">
                <h2 class="font-cursive text-5xl text-gold text-center mb-16">Love Story</h2>
                
                <div class="relative">
                    <div class="timeline-line hidden md:block"></div>
                    
                    <!-- Story Items -->
                    <div class="space-y-12">
                        <div class="relative flex flex-col md:flex-row items-center gap-8">
                            <div class="flex-1 md:text-right">
                                <h3 class="font-cursive text-2xl text-gold mb-2">Pertemuan Pertama</h3>
                                <p class="text-sm text-gray-500 mb-2">Januari 2020</p>
                                <p class="text-gray-600">Pertemuan pertama yang tak terduga, mengawali kisah indah kami berdua.</p>
                            </div>
                            <div class="w-4 h-4 bg-gold rounded-full border-4 border-cream shadow-lg z-10"></div>
                            <div class="flex-1"></div>
                        </div>
                        
                        <div class="relative flex flex-col md:flex-row items-center gap-8">
                            <div class="flex-1"></div>
                            <div class="w-4 h-4 bg-gold rounded-full border-4 border-cream shadow-lg z-10"></div>
                            <div class="flex-1">
                                <h3 class="font-cursive text-2xl text-gold mb-2">Memulai Hubungan</h3>
                                <p class="text-sm text-gray-500 mb-2">Juni 2020</p>
                                <p class="text-gray-600">Memutuskan untuk berjalan bersama dalam satu visi dan misi.</p>
                            </div>
                        </div>
                        
                        <div class="relative flex flex-col md:flex-row items-center gap-8">
                            <div class="flex-1 md:text-right">
                                <h3 class="font-cursive text-2xl text-gold mb-2">Lamaran</h3>
                                <p class="text-sm text-gray-500 mb-2">Desember 2024</p>
                                <p class="text-gray-600">Dengan restu kedua orang tua, kami berkomitmen untuk melangkah ke jenjang pernikahan.</p>
                            </div>
                            <div class="w-4 h-4 bg-gold rounded-full border-4 border-cream shadow-lg z-10"></div>
                            <div class="flex-1"></div>
                        </div>
                    </div>
                </div>
            </div>
        </section>

        <!-- Gallery Section -->
        <section class="py-24 px-6 bg-primary/30">
            <h2 class="font-cursive text-5xl text-gold text-center mb-16">Our Gallery</h2>
            <div class="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4 max-w-6xl mx-auto">
                {{range .Gallery}}
                <div class="aspect-square overflow-hidden rounded-xl shadow-lg hover:shadow-2xl transition-all duration-500 cursor-pointer group">
                    <img src="{{.URL}}" class="w-full h-full object-cover group-hover:scale-110 transition-transform duration-700" alt="Gallery">
                </div>
                {{end}}
            </div>
        </section>

        <!-- Wedding Gift Section -->
        <section class="py-24 px-6 bg-cream">
            <div class="max-w-xl mx-auto text-center">
                <h2 class="font-cursive text-5xl text-gold mb-4">Wedding Gift</h2>
                <p class="text-gray-600 mb-12">Doa restu Anda adalah hadiah terindah bagi kami. Namun jika Anda ingin memberikan tanda kasih, kami menyediakan amplop digital.</p>
                
                <div class="space-y-6">
                    <!-- Bank BCA -->
                    <div class="bank-bca rounded-2xl p-6 text-white text-left relative overflow-hidden">
                        <div class="absolute top-4 right-4 w-12 h-12 bg-white/20 rounded-full"></div>
                        <p class="text-sm opacity-80 mb-4">Bank BCA</p>
                        <p class="text-xl font-bold mb-1">{{.BrideName}}</p>
                        <p class="text-2xl font-mono tracking-wider mb-4">1234567890</p>
                        <button onclick="copyToClipboard('1234567890')" class="flex items-center gap-2 text-sm bg-white/20 hover:bg-white/30 px-4 py-2 rounded-lg transition">
                            <svg class="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 16H6a2 2 0 01-2-2V6a2 2 0 012-2h8a2 2 0 012 2v2m-6 12h8a2 2 0 002-2v-8a2 2 0 00-2-2h-8a2 2 0 00-2 2v8a2 2 0 002 2z"/>
                            </svg>
                            Salin Nomor
                        </button>
                    </div>
                </div>
            </div>
        </section>

        <!-- RSVP Section -->
        <section class="py-24 px-6 bg-primary/30">
            <div class="max-w-xl mx-auto">
                <h2 class="font-cursive text-5xl text-gold text-center mb-4">RSVP</h2>
                <p class="text-gray-600 text-center mb-12">Konfirmasi kehadiran dan kirimkan ucapan untuk kedua mempelai</p>
                
                <!-- HTMX RSVP Form -->
                <div hx-get="/invitations/{{.Subdomain}}/rsvp" hx-trigger="load"></div>
                
                <!-- HTMX Messages -->
                <div class="mt-12" hx-get="/invitations/{{.Subdomain}}/messages" hx-trigger="load"></div>
            </div>
        </section>

        <!-- Footer -->
        <section class="py-16 px-6 cover-bg relative">
            <div class="absolute inset-0 bg-black/70"></div>
            <div class="relative z-10 text-center text-white">
                <p class="text-lg mb-4">Terima kasih atas doa dan restu yang diberikan</p>
                <h2 class="font-cursive text-5xl text-gold mb-4">{{.GroomNickname}} & {{.BrideNickname}}</h2>
                <p class="text-sm opacity-70">&copy; 2025 - Created with Love</p>
            </div>
        </section>
    </main>

    <script>
        // Open Invitation
        function openInvitation() {
            document.getElementById('opening').classList.add('hidden');
            document.getElementById('content').classList.remove('hidden');
            
            // Auto play music
            const audio = document.getElementById('bgm');
            if(audio) {
                audio.play().catch(e => console.log("Audio autoplay blocked", e));
            }
            
            // Start countdown
            startCountdown();
        }
        
        // Toggle Music
        let isPlaying = true;
        function toggleMusic() {
            const audio = document.getElementById('bgm');
            if (!audio) return;
            
            if (audio.paused) {
                audio.play();
                isPlaying = true;
            } else {
                audio.pause();
                isPlaying = false;
            }
        }
        
        // Countdown Timer
        function startCountdown() {
            const weddingDate = new Date('{{.WeddingDateISO}}').getTime();
            
            const timer = setInterval(function() {
                const now = new Date().getTime();
                const distance = weddingDate - now;
                
                if (distance < 0) {
                    clearInterval(timer);
                    document.getElementById('days').textContent = '00';
                    document.getElementById('hours').textContent = '00';
                    document.getElementById('minutes').textContent = '00';
                    document.getElementById('seconds').textContent = '00';
                    return;
                }
                
                const days = Math.floor(distance / (1000 * 60 * 60 * 24));
                const hours = Math.floor((distance % (1000 * 60 * 60 * 24)) / (1000 * 60 * 60));
                const minutes = Math.floor((distance % (1000 * 60 * 60)) / (1000 * 60));
                const seconds = Math.floor((distance % (1000 * 60)) / 1000);
                
                document.getElementById('days').textContent = String(days).padStart(2, '0');
                document.getElementById('hours').textContent = String(hours).padStart(2, '0');
                document.getElementById('minutes').textContent = String(minutes).padStart(2, '0');
                document.getElementById('seconds').textContent = String(seconds).padStart(2, '0');
            }, 1000);
        }
        
        // Copy to Clipboard
        function copyToClipboard(text) {
            navigator.clipboard.writeText(text).then(() => {
                alert('Nomor rekening berhasil disalin!');
            }).catch(err => {
                console.error('Failed to copy: ', err);
            });
        }
    </script>
</body>
</html>`,
			DefaultCSS: `/* Luxury Premium Theme CSS */
html {
    scroll-behavior: smooth;
}
body {
    -webkit-font-smoothing: antialiased;
    -moz-osx-font-smoothing: grayscale;
}
/* Custom scrollbar */
::-webkit-scrollbar {
    width: 8px;
}
::-webkit-scrollbar-track {
    background: #f1f1f1;
}
::-webkit-scrollbar-thumb {
    background: #C5A880;
    border-radius: 4px;
}`,
		},
	}

	for _, cat := range categories {
		var existing theme.ThemeCategory
		err := db.Where("slug = ?", cat.Slug).First(&existing).Error
		if err != nil {
			db.Create(&cat)
			fmt.Printf("Created category: %s\n", cat.Name)
		}
	}
}
