<script setup lang="ts">
import HelloWorld from './components/HelloWorld.vue'

import {GrpcWebFetchTransport} from "@protobuf-ts/grpcweb-transport";
import {MyServiceClient} from "./backend/gRPC/service.client";
import {HelloRequest, TimeRequest, TimeResponse} from "./backend/gRPC/service";

// Create a transport
const transport = new GrpcWebFetchTransport({
  baseUrl: "http://localhost:8080", // The address of your gRPC-Web proxy
  withCredentials: false,
});

// Create a client
const client = new MyServiceClient(transport);

// Create a request
const request: HelloRequest = {
  name: "World",
};

// Call the gRPC method
client.sayHello(request).then(response => {
  console.log(response.response.message);
}).catch(error => {
  console.error("Error calling gRPC method:", error);
});

// Ensure TimeRequest is instantiated correctly
const requestTime: TimeRequest = {};
const stream = client.streamTime(requestTime);

// Subscribe to the stream
stream.responses.onNext((response: TimeResponse) => {
  console.log('Received response:', response.currentTime);
});

stream.responses.onError(error => {
  console.error('Stream error:', error);
  currentTime.value = 'Error occurred while fetching time';
});

stream.responses.onComplete(() => {
  console.log('Stream completed');
});

</script>

<template>
  <div>
    <a href="https://vitejs.dev" target="_blank">
      <img src="/vite.svg" class="logo" alt="Vite logo"/>
    </a>
    <a href="https://vuejs.org/" target="_blank">
      <img src="./assets/vue.svg" class="logo vue" alt="Vue logo"/>
    </a>
  </div>
  <HelloWorld msg="Vite + Vue"/>
</template>

<style scoped>
.logo {
  height: 6em;
  padding: 1.5em;
  will-change: filter;
  transition: filter 300ms;
}

.logo:hover {
  filter: drop-shadow(0 0 2em #646cffaa);
}

.logo.vue:hover {
  filter: drop-shadow(0 0 2em #42b883aa);
}
</style>
