document.addEventListener('DOMContentLoaded', () => {
    const newsContainer = document.getElementById('news-container');
    const categoryFilters = document.getElementById('category-filters');
    const userInfo = document.getElementById('user-info');
    const userName = document.getElementById('user-name');
    const userPicture = document.getElementById('user-picture');

    let newsData = [];
    let selectedCategories = [];

    // ফেক ইউজার ডেটা (সাধারণত এটি সার্ভার থেকে আসবে)
    const user = JSON.parse(localStorage.getItem('user')); // লগইন করা ইউজারের তথ্য

    if (user) {
        userName.textContent = user.name; // ইউজারের নাম দেখানো
        userPicture.src = user.picture; // ইউজারের ছবি দেখানো
        userInfo.style.display = 'flex'; // ইউজার তথ্য দেখানো
    }

    async function loadNews() {
        try {
            const response = await fetch('http://localhost:8080/news');
            if (!response.ok) {
                throw new Error('Failed to fetch news');
            }

            newsData = await response.json();
            renderNews(newsData);
            renderCategoryFilters(newsData);
        } catch (error) {
            console.error(error);
            newsContainer.innerHTML = '<p>Error loading news.</p>';
        }
    }

    function renderNews(filteredNews) {
        newsContainer.innerHTML = ''; // পূর্ববর্তী নিউজ ক্লিয়ার করা

        if (filteredNews.length === 0) {
            newsContainer.innerHTML = '<p>No news available for the selected category.</p>';
            return;
        }

        filteredNews.forEach(news => {
            const newsItem = document.createElement('div');
            newsItem.classList.add('news-item');

            newsItem.innerHTML = `
                <h2>${news.title}</h2>
                <img src="${news.image}" alt="${news.title}" onerror="this.onerror=null; this.src='fallback-image.jpg';">
                <p>${news.description}</p>
                <p><strong>Category:</strong> ${news.category}</p>
                <p><strong>Date:</strong> ${news.date}</p>
            `;

            newsContainer.appendChild(newsItem);
        });
    }

    function renderCategoryFilters(newsData) {
        const categories = [...new Set(newsData.map(news => news.category))]; // ইউনিক ক্যাটাগরি এক্সট্র্যাক্ট করা
        categoryFilters.innerHTML = ''; // পূর্ববর্তী ফিল্টার ক্লিয়ার করা

        categories.forEach(category => {
            const label = document.createElement('label');
            label.innerHTML = `
                <input type="checkbox" value="${category}"> ${category}
            `;
            categoryFilters.appendChild(label);
        });

        categoryFilters.addEventListener('change', handleCategoryFilter);
    }

    function handleCategoryFilter() {
        selectedCategories = Array.from(categoryFilters.querySelectorAll('input[type="checkbox"]:checked')).map(checkbox => checkbox.value);
        
        if (selectedCategories.length === 0) {
            renderNews(newsData); // কোনো ক্যাটাগরি নির্বাচিত না হলে সমস্ত নিউজ দেখানো
        } else {
            const filteredNews = newsData.filter(news => selectedCategories.includes(news.category));
            renderNews(filteredNews);
        }
    }

    loadNews();
});
