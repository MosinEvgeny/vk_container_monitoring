<template>
  <div class="container">
    <h1>Container Monitoring</h1>
    <table class="table">
      <thead>
        <tr>
          <th>IP Address</th>
          <th>Ping Time (ms)</th>
          <th>Last Successful Ping</th>
        </tr>
      </thead>
      <tbody>
        <tr v-for="ping in pings" :key="ping.ip_address">
          <td>{{ ping.ip_address }}</td>
          <td>{{ ping.ping_time }}</td>
          <td>{{ ping.last_successful_ping }}</td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

<script>
import { ref, onMounted } from 'vue';
import 'bootstrap/dist/css/bootstrap.min.css';
import axios from 'axios';

export default {
  setup() {
    const pings = ref([]);

    const fetchData = async () => {
      try {
        const response = await axios.get(import.meta.env.VITE_BACKEND_URL + '/pings');
        pings.value = response.data;
      } catch (error) {
        console.error("Could not fetch pings:", error);
      }
    };

    onMounted(() => {
      fetchData();
    });

    return {
      pings,
    };
  }
};
</script>

<style scoped></style>