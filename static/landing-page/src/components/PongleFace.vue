<template>
  <div class="pongle">
    <img
      class="face"
      alt="Vue logo"
      src="../assets/pongle.png"
      width="128"
      height="128"
    />
    <img
      class="happy"
      alt="Vue logo"
      src="../assets/happy.png"
      width="128"
      height="128"
    />
    <img
      class="paw"
      alt="Vue logo"
      src="../assets/paw.png"
      width="128"
      height="128"
    />
  </div>
</template>

<script>
export default {
  name: "HelloWorld",
  props: {
    msg: String,
  },
};
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
.pongle {
  display: grid;
  grid-template-columns: auto 128px 128px 128px auto;
}

.face {
  grid-column: 3;
  grid-row: 1;
  justify-self: center;
  animation-duration: 5s;
  animation-name: happy-out;
  animation-iteration-count: infinite;
  animation-delay: 5s;
}

.happy {
  grid-column: 3;
  grid-row: 1;
  /* opacity: 0; */
  justify-self: center;
  /* animation-duration: 5s;
  animation-name: happy-in;
  animation-iteration-count: infinite;
  animation-delay: 5s; */
}

.paw {
  grid-column: 4;
  justify-self: left;
  animation-duration: 5s;
  animation-name: wave;
  animation-iteration-count: infinite;
  animation-delay: 5s;
  transform: translate(0, 50px) rotateX(90deg) translate(0, -32px);
}

@keyframes happy-in {
  from {
    opacity: 0;
  }

  35% {
    opacity: 0;
  }

  44% {
    opacity: 1;
  }

  56% {
    opacity: 1;
  }

  65% {
    opacity: 0;
  }

  to {
    opacity: 0;
  }
}

@keyframes happy-out {
  from {
    opacity: 1;
  }

  35% {
    opacity: 1;
  }

  44% {
    opacity: 0;
  }

  56% {
    opacity: 0;
  }

  65% {
    opacity: 1;
  }

  to {
    opacity: 1;
  }
}

@keyframes wave {
  from {
    transform: translate(0, 50px) rotateX(90deg) translate(0, -32px);
  }

  30% {
    transform: translate(0, 50px) rotateX(90deg) translate(0, -32px);
  }

  35% {
    transform: translate(0, -30px) rotateX(0deg) translate(0, -32px);
  }

  44% {
    transform: translate(0, -30px) rotateZ(0deg) translate(0, -32px);
  }

  46% {
    transform: translate(0, -30px) rotateZ(10deg) translate(0, -32px);
  }

  48% {
    transform: translate(0, -30px) rotateZ(-10deg) translate(0, -32px);
  }

  50% {
    transform: translate(0, -30px) rotateZ(10deg) translate(0, -32px);
  }

  52% {
    transform: translate(0, -30px) rotateZ(-10deg) translate(0, -32px);
  }

  54% {
    transform: translate(0, -30px) rotateZ(10deg) translate(0, -32px);
  }

  56% {
    transform: translate(0, -30px) rotateZ(0deg) translate(0, -32px);
  }

  61% {
    transform: translate(0, 50px) rotateX(90deg) translate(0, -32px);
  }

  to {
    transform: translate(0, 50px) rotateX(90deg) translate(0, -32px);
  }
}
</style>
