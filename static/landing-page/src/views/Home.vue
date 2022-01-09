<template>
  <div class="home">
    <HelloWorld msg="Welcome to PongleHub Games!" />
    <PongleFace />
  </div>
</template>

<script>
// @ is an alias to /src
import HelloWorld from "@/components/HelloWorld.vue";
import PongleFace from "@/components/PongleFace.vue";

export default {
  name: "Home",
  components: {
    HelloWorld,
    PongleFace,
  },
};
</script>
