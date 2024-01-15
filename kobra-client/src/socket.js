import { reactive } from 'vue';
import { io } from 'socket.io-client';

// Create a reactive state to store the Socket.IO connection status and events
export const state = reactive({
  connected: false,
  events: []
});

// Create a Socket.IO connection to the Python Socket.IO service
const socket = io("http://127.0.0.1:5000"); // Replace with your Python Socket.IO service URL

// Listen for the 'connect' event
socket.on('connect', () => {
  state.connected = true;
});

// Listen for the 'disconnect' event
socket.on('disconnect', () => {
  state.connected = false;
});

// Listen for a custom event, replace 'event_name' with your actual event name
socket.on('event_name', (data) => {
  state.events.push(data);
});



// Export the socket so you can use it to emit events to the server
export default socket;