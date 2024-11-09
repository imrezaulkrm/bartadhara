document.addEventListener('DOMContentLoaded', () => {
    const newsContainer = document.getElementById('news-container');
    const categoryFilters = document.getElementById('category-filters');
    const userInfo = document.getElementById('user-info');
    const userName = document.getElementById('user-name');
    const userPicture = document.getElementById('user-picture');

    let newsData = [];
    let selectedCategories = [];

    // Fake user data (Normally this would come from the server)
    const user = JSON.parse(localStorage.getItem('user')); // Logged-in user info

    if (user) {
        userName.textContent = user.name; // Show the user's name
        userPicture.src = user.picture; // Show the user's picture
        userInfo.style.display = 'flex'; // Display user info
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
        newsContainer.innerHTML = ''; // Clear previous news

        if (filteredNews.length === 0) {
            newsContainer.innerHTML = '<p>No news available for the selected category.</p>';
            return;
        }

        filteredNews.forEach(news => {
            const newsItem = document.createElement('div');
            newsItem.classList.add('news-item');

            newsItem.innerHTML = `
                <img src="${news.image}" alt="${news.title}" onerror="this.onerror=null; this.src='fallback-image.jpg';">
                <h2>${news.title}</h2>
                <p class="description">${news.description}</p>
                <span class="see-more" onclick="viewNewsDetails('${news.id}')">See more</span>
                <div class="footer">
                    <span>${news.date}</span>
                    <span>${news.category}</span>
                </div>
            `;

            newsContainer.appendChild(newsItem);
        });
    }

    function renderCategoryFilters(newsData) {
        const categories = [...new Set(newsData.map(news => news.category))]; // Extract unique categories
        categoryFilters.innerHTML = ''; // Clear previous filters

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
            renderNews(newsData); // Show all news if no category is selected
        } else {
            const filteredNews = newsData.filter(news => selectedCategories.includes(news.category));
            renderNews(filteredNews);
        }
    }

    function viewNewsDetails(newsID) {
        console.log("Navigating to news details with ID:", newsID); // Log newsID for debugging
        window.location.href = `news-details.html?id=${newsID}`; // Navigate to the news details page
    }

    loadNews();
});
