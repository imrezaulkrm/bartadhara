document.addEventListener('DOMContentLoaded', async () => {
    // URL থেকে 'id' প্যারামিটার নেয়া
    const urlParams = new URLSearchParams(window.location.search);
    const newsID = urlParams.get('id');

    if (newsID) {
        try {
            const response = await fetch(`http://localhost:8080/news/${newsID}`);
            if (!response.ok) {
                throw new Error('Failed to fetch news details');
            }

            const news = await response.json();

            // নিউজ ডেটা এই HTML এলিমেন্টে সেট করা
            document.getElementById('news-title').textContent = news.title;
            document.getElementById('news-image').src = news.image;
            document.getElementById('news-description').textContent = news.description;
            document.getElementById('news-category').textContent = `Category: ${news.category}`;
            document.getElementById('news-date').textContent = `Published on: ${news.date}`;

        } catch (error) {
            console.error(error);
            document.getElementById('news-details').innerHTML = '<p>Error loading news details.</p>';
        }
    } else {
        console.error('No news ID found in URL');
        document.getElementById('news-details').innerHTML = '<p>No news found.</p>';
    }
});
