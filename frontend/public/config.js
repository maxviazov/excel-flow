// API Configuration
const API_BASE_URL = window.location.hostname === 'localhost' 
    ? 'http://localhost:8080'
    : window.location.protocol === 'https:'
        ? 'https://excel-flow-alb-1086104942.us-east-1.elb.amazonaws.com'
        : 'http://excel-flow-alb-1086104942.us-east-1.elb.amazonaws.com';
