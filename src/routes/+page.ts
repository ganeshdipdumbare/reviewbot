const backendUrl = import.meta.env.VITE_BACKEND_URL;

export const load = async ({ fetch }) => {
	console.log('fetching initial message');
	const response = await fetch(`${backendUrl}/api/v1/converse`, {
		method: 'POST',
		headers: {
			'Content-Type': 'application/json'
		},
		body: JSON.stringify({ text: 'system@initiateconversation' })
	});
	console.log(response);
	if (!response.ok) {
		throw new Error('Failed to fetch initial message');
	}

	const data = await response.json();
	return {
		sender: 'user',
		conversationID: data.conversationID,
		reviewID: data.reviewID,
		userID: 'user1',
		productID: 'product1',
		text: data.text
	};
};
