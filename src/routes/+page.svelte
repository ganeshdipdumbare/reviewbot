<script lang="ts">
	import { onMount, afterUpdate } from 'svelte';
	import type { PageData } from './$types';
	export let data: PageData;
	const imageUrl = '/images/bot.png';
	const backendUrl = import.meta.env.VITE_BACKEND_URL;
	interface ChatMessage {
		sender: 'user' | 'bot';
		conversationID?: string;
		reviewID?: string;
		userID?: string;
		productID?: string;
		text: string;
	}

	let requestBody: ChatMessage = {
		sender: 'user',
		conversationID: '',
		reviewID: '',
		userID: 'user1',
		productID: 'product1',
		text: ''
	};
	let chatMessages: ChatMessage[] = [];
	let error: string | null = null;
	let initialMessage = '';

	const sendMessage = async () => {
		try {
			console.log('Sending message:', requestBody);
			const response = await fetch(`${backendUrl}/api/v1/converse`, {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json'
				},
				body: JSON.stringify(requestBody)
			});

			if (!response.ok) {
				throw new Error('Failed to send message');
			}

			const data = await response.json();
			chatMessages = [
				...chatMessages,
				{ sender: 'user', text: requestBody.text },
				{ sender: 'bot', text: data.text }
			];

			// Update requestBody for the next request
			requestBody.conversationID = data.conversationID;
			requestBody.reviewID = data.reviewID;
			requestBody.text = ''; // Clear the input text
			error = null;
		} catch (err) {
			error = err instanceof Error ? err.message : 'An unknown error occurred';
		}
	};

	const systemInitiateReview = async () => {
		try {
			requestBody.text = 'system@startreview';
			console.log('Sending message:', requestBody);
			const response = await fetch(`${backendUrl}/api/v1/converse`, {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json'
				},
				body: JSON.stringify(requestBody)
			});

			if (!response.ok) {
				throw new Error('Failed to send message');
			}

			const data = await response.json();
			chatMessages = [
				...chatMessages,
				//{ sender: 'user', text: requestBody.text },
				{ sender: 'bot', text: data.text }
			];

			// Update requestBody for the next request
			requestBody.conversationID = data.conversationID;
			requestBody.reviewID = data.reviewID;
			requestBody.text = ''; // Clear the input text
			error = null;
		} catch (err) {
			error = err instanceof Error ? err.message : 'An unknown error occurred';
		}
	};

	const endConversation = async () => {
		try {
			requestBody.text = 'bye';
			console.log('Sending message:', requestBody);
			const response = await fetch(`${backendUrl}/api/v1/converse`, {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json'
				},
				body: JSON.stringify(requestBody)
			});

			if (!response.ok) {
				throw new Error('Failed to send message');
			}

			const data = await response.json();
			chatMessages = [
				...chatMessages,
				//{ sender: 'user', text: requestBody.text },
				{ sender: 'bot', text: data.text }
			];

			// Update requestBody for the next request
			requestBody.conversationID = '';
			requestBody.reviewID = '';
			requestBody.text = '';
			requestBody.userID = '';
			requestBody.productID = '';
			error = null;
		} catch (err) {
			error = err instanceof Error ? err.message : 'An unknown error occurred';
		}
	};

	const scrollToBottom = () => {
		const container = document.querySelector('.message-container');
		if (container) {
			container.scrollTop = container.scrollHeight;
		}
	};

	onMount(() => {
		// Set initialMessage from the page store data
		initialMessage = data.text;
		chatMessages = [{ sender: 'bot', text: initialMessage }];
		// Set the initial message in the requestBody
		requestBody.conversationID = data.conversationID;
		requestBody.reviewID = data.reviewID;
		requestBody.userID = data.userID;
		requestBody.productID = data.productID;
		scrollToBottom();
	});

	afterUpdate(() => {
		scrollToBottom();
	});

	const handleKeyDown = (event: KeyboardEvent) => {
		if (event.key === 'Enter') {
			event.preventDefault(); // Prevent default form submission behavior
			sendMessage();
		}
	};
</script>

<div class="chat-container">
	<div class="header">
		<img src={imageUrl} alt="Bot Icon" class="bot-icon" style="width: 100px; height: 100px;" />
		<h1 class="bot-name">ReviewBot</h1>
	</div>

	<div class="message-container">
		{#each chatMessages as { sender, text }, index (index)}
			<div class="message {sender}">
				{text}
			</div>
		{/each}
	</div>

	<div class="input-container">
		<input
			type="text"
			bind:value={requestBody.text}
			placeholder="Type your message..."
			on:keydown={handleKeyDown}
		/>
		<button on:click={sendMessage}>Send</button>
	</div>
	<div class="button-container">
		<button class="gray-button" on:click={systemInitiateReview}>Order delivered</button>
		<button class="red-button" on:click={endConversation}>End conversation</button>
	</div>
	{#if error}
		<p style="color: red;">{error}</p>
	{/if}
</div>

<style>
	.chat-container {
		width: 100%;
		max-width: 600px;
		margin: auto;
		padding: 20px;
		border: 1px solid #ccc;
		border-radius: 8px;
		display: flex;
		flex-direction: column;
		gap: 10px;
		height: 80vh; /* Adjust height as needed */
	}

	.message-container {
		display: flex;
		flex-direction: column;
		gap: 10px;
		flex-grow: 1;
		overflow-y: auto;
	}

	.input-container {
		display: flex;
		gap: 10px;
	}

	.message {
		margin: 0;
		padding: 10px;
		border-radius: 8px;
		max-width: fit-content;
		word-wrap: break-word;
	}

	.user {
		background-color: #e6f7ff;
		align-self: flex-end;
	}

	.bot {
		background-color: #f0f0f0;
		align-self: flex-start;
	}

	input[type='text'] {
		flex-grow: 1;
		padding: 10px;
		border: 1px solid #ccc;
		border-radius: 8px;
	}

	button {
		padding: 10px 20px;
		border: none;
		border-radius: 8px;
		background-color: #007bff;
		color: white;
		cursor: pointer;
	}

	.button-container {
		display: flex;
		justify-content: space-between;
		margin-top: 10px;
	}

	.gray-button {
		background-color: gray;
	}

	.red-button {
		background-color: red;
	}

	button:hover {
		background-color: #0056b3;
	}
</style>
