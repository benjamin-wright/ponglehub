<template>
  <div class="hello">
    <h1>{{ msg }}</h1>
    <p>{{ done }}</p>
  </div>
</template>

<script>
import axios from 'axios'

axios.defaults.withCredentials = true

export default {
  name: 'HelloWorld',
  props: {
    msg: String
  },
  data() {
    return {
      result: "",
      done: false
    }
  },
  mounted() {
    axios
      .get('https://game-state.ponglehub.co.uk/status')
      .then(response => {
        this.result = response.data
      })
      .catch(error => this.result = error.toString())
      .finally(() => this.done = true)

  }
}
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
h3 {
  margin: 40px 0 0;
}
ul {
  list-style-type: none;
  padding: 0;
}
li {
  display: inline-block;
  margin: 0 10px;
}
a {
  color: #42b983;
}
</style>
