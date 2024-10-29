document.getElementById('registration-form').addEventListener('submit', async (e) => {
    e.preventDefault(); // ফর্মের ডিফল্ট সাবমিশন আচরণ বন্ধ করা

    // ফর্মের ইনপুট থেকে মান সংগ্রহ করা
    const name = document.getElementById('name').value;
    const username = document.getElementById('username').value;
    const email = document.getElementById('email').value;
    const password = document.getElementById('password').value;
    const pictureInput = document.getElementById('picture');

    // FormData অবজেক্ট তৈরি করা
    const formData = new FormData();
    formData.append('name', name);
    formData.append('username', username);
    formData.append('email', email);
    formData.append('password', password);

    // ছবি ফাইল চেক করা এবং FormData-এ যোগ করা
    if (pictureInput.files.length > 0) {
        formData.append('picture', pictureInput.files[0]); // ছবির ফাইল যোগ করা
    }

    try {
        // সার্ভারে ডেটা পাঠানো
        const response = await fetch('http://localhost:8080/users', {
            method: 'POST',
            body: formData, // FormData অবজেক্ট পাঠানো
        });

        // সার্ভার থেকে পাওয়া প্রতিক্রিয়া চেক করা
        if (response.ok) {
            alert('Registration successful!'); // সফল হলে বার্তা
            window.location.href = 'login.html'; // লগইন পেইজে রিডাইরেক্ট
        } else {
            const errorData = await response.json(); // ত্রুটির ডেটা গ্রহণ করা
            alert(`Registration failed: ${errorData.message || 'Please try again.'}`); // ব্যর্থ হলে বার্তা
        }
    } catch (error) {
        console.error('Error:', error); // কনসোলে ত্রুটি লগ করা
        alert('An unexpected error occurred. Please try again.'); // অপ্রত্যাশিত ত্রুটি হলে বার্তা
    }
});
