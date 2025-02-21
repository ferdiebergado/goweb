const form = document.getElementById('x-form') as HTMLFormElement;

form.addEventListener('submit', async (e) => {
  e.preventDefault();
  try {
    const frmData = new FormData(form);
    const obj = Object.fromEntries(frmData);
    const opts: RequestInit = {
      method: form.method,
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(obj),
    };
    const res = await fetch(form.action, opts);
    if (res.ok) {
      const data = await res.json();
      console.log(data);
    }
  } catch (error) {
    console.error(error);
  }
});
