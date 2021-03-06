const forms = ['signup', 'login']
const serverAddress = 'https://socialnet.rodolforg.com/api'

const addClickListener = form => {
  document.querySelector(`#${form} .submit`).addEventListener('click', async () => {
    let postData = {}
    postData.username = document.querySelector(`#${form} .username`).value
    postData.password = document.querySelector(`#${form} .password`).value
    if (form === 'signup') {
      postData.firstName = document.querySelector(`#${form} .first-name`).value
      postData.lastName = document.querySelector(`#${form} .last-name`).value
      postData.email = document.querySelector(`#${form} .email`).value
    }

    let token
    try {
      token = await fetch(`${serverAddress}/${form}`, {
        method: 'post',
        headers: new Headers({
          'Content-Type': 'Application/json'
        }),
        body: JSON.stringify(postData)
      })
      .catch(console.error)
      .then(res => res.json())
      .then(json => json.token)
    } catch (err) {
      console.error(err)
      return
    }

    console.log('token: ', token)
    document.cookie = `socialnet_token=${token};`
    window.location = `/user/${postData.username}`
  })
}
  
forms.forEach(addClickListener)
