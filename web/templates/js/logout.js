function Logout(home){
	console.log("i got called Logout .... ");
	axios.get(
		"/api/logout",
		{
			responseType:"application/json"
		}
	).then(function(response){
		var data =response.data;
		console.log(data);
		if (data.success){
			var elem = document.createElement("a");

			// var el = document.getElementById("")

			elem.setAttribute("href"  , home );
			elem.click();
		}else {
			console.log("can't Log in nigga ");
		}
	});
}