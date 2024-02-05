const  express = require('express');
const router = express.Router();

const app = express()


app.set('port',process.env.PORT || 3200);

app.use(express.json());

app.use('/', router)

router.post("/", async (req, res)=>{
    const {name, email,  comment} = req.body;
    
    if(!name || !email || !comment){
        return res.status(400).send("Missing fields")   
        }
        
       let data = await sendEmail(name, email , comment);
       console.log(data)
       res.json({msg: "Message sent"})
   })

app.listen(app.get('port'), () => console.log(`listening on:${app.get('port')}`));
