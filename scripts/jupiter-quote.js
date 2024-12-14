import http from 'k6/http';

export const options = {
    vus: 500,
    duration: "10s",
};

export default function () {
    const response = http.get('http://0.0.0.0:8081/quote?inputMint=So11111111111111111111111111111111111111112&outputMint=DezXAZ8z7PnrnRJjz3wXBoRgixCa6xjnB7YaB1pPB263&amount=1000000');
}
