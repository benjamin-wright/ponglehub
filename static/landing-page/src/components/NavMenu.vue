<template>
  <a class="button" v-on:click="toggle()">...</a>
  <transition name="slide" :duration="500">
    <div v-if="expanded" class="expander">
      <a v-on:click="this.toggle()">...</a>
      <a
        v-for="item in items"
        :key="item"
        v-on:click="
          this.$emit('select:option', item);
          this.toggle();
        "
      >
        {{ item }}
      </a>
      <a
        v-on:click="
          this.$emit('select:logout');
          this.toggle();
        "
      >
        logout
      </a>
    </div>
  </transition>
</template>

<script>
export default {
  name: "NavMenu",
  props: ["items"],
  emits: ["select:option", "select:logout"],
  data: () => {
    return {
      expanded: false,
    };
  },
  methods: {
    toggle() {
      this.expanded = !this.expanded;
    },
  },
};
</script>

<!-- Add "scoped" attribute to limit CSS to this component only -->
<style scoped>
.button {
  cursor: pointer;
  user-select: none;
  padding: 0 0.5em;
}

.expander {
  position: absolute;
  top: 0;
  right: 0;

  background: #2c3e50;
  transform: translateZ(-10px);

  display: flex;
  flex-direction: column;
  text-transform: capitalize;
  align-items: flex-end;

  border-radius: 0 0 0 0.5em;
}

.slide-enter-active, .slide-leave-active {
  transition: transform 0.5s;
}

.slide-enter-from, .slide-leave-to {
  transform: translateZ(-10px) translateY(-100%);
}

.slide-enter-to, .slide-leave-from {
  transform: translateZ(-10px) translateY(0%);
}

.expander a {
  padding: 0.5em;
  margin: 0.5em;
  margin-right: 1em;

  cursor: pointer;
  user-select: none;
}
</style>
