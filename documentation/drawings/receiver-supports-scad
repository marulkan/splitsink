$fn=60;
// Audio backplate
module backplate(){
  cube([80,7,5]);
  translate([3.5,3.5,0]) {
    cylinder(h=10, r=1.25);
    translate([73,0,0]) {
      cylinder(h=10, r=1.25);
    }
  }
  translate([0,103,0]){
    cube([80,7,5]);
    translate([3.5,3.5,0]) {
      cylinder(h=10, r=1.25);
      translate([73,0,0]) {
        cylinder(h=10, r=1.25);
      }
    }
  }
  cube([7,110,5]);
  translate([73,0,0]) 
  cube([7,110,5]);

  translate([-21,0,0]) {
    difference(){
    cube([21,43,5]);
    translate([3.5,3.5,0])
    cube([21-7,43-7,5]);
    }
    translate([2,6,0])
    cylinder(h=10, r=1.25);

    translate([21-2,43-6,0])
    cylinder(h=10, r=1.25);
    translate([-2,0,0])
    // padding
    cube([2,43,5]);
  }
}

backplate();
translate([80+21+2,0,0])
backplate();
translate([2*(80+21+2),0,0])
backplate();
