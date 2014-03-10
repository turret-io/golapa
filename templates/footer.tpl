<div class="footer-cont">
    <div class="footer-inner">
        <div class="footer-left">
            <p class="lead">
                <a href="mailto:hello@yourdomain.com">hello@yourdomain.com</a>
            </p>
        </div>
        <div class="footer-copy">
            <p>&copy; <script>document.write(new Date().getFullYear())</script>&nbsp;yourdomain.com</p>
        </div>
    </div>
</div>


             <script src="//ajax.googleapis.com/ajax/libs/jquery/1.10.1/jquery.min.js"></script>
        <script>window.jQuery || document.write('<script src="js/vendor/jquery-1.10.1.min.js"><\/script>')</script>

        <script src="/js/vendor/bootstrap.min.js"></script>

        <script src="/js/main.js"></script>

        
    <div id="fb-root"></div>

<script type="text/javascript">
    jQuery(document).ready(function(){
        jQuery(".submit-signup").off('submit').on('submit', function(e){
            var errors = Array();
            jQuery(e.target).find(":input").each(function(k,v) {
                if(jQuery(v).val().length < 1) {
                    if(jQuery(v)[0].type != "submit") {
                        errors.push(v.name);
                        jQuery(v).parents(".form-group").addClass("has-error");
                    }
                }
            });
            if(errors.length > 0) {
                e.preventDefault();
            }

        });
    });
</script>
</body>
</html>